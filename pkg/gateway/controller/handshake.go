// Copyright 2023 Ant Group Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/secretflow/kuscia/pkg/common"
	kusciaapisv1alpha1 "github.com/secretflow/kuscia/pkg/crd/apis/kuscia/v1alpha1"
	"github.com/secretflow/kuscia/pkg/gateway/clusters"
	"github.com/secretflow/kuscia/pkg/gateway/config"
	"github.com/secretflow/kuscia/pkg/gateway/xds"
	"github.com/secretflow/kuscia/pkg/utils/nlog"
	tlsutils "github.com/secretflow/kuscia/pkg/utils/tls"
	"github.com/secretflow/kuscia/proto/api/v1alpha1"
	"github.com/secretflow/kuscia/proto/api/v1alpha1/handshake"
)

const (
	tokenByteSize = 32
	NoopToken     = "noop"
)

const (
	handShakeTypeUID = "UID"
	handShakeTypeRSA = "RSA"
)

var (
	tokenPrefix = []byte("kuscia")
)

type Token struct {
	Token   string
	Version int64
}

type AfterRegisterDomainHook func(response *handshake.RegisterResponse)

func (c *DomainRouteController) startHandShakeServer(port uint32) {
	mux := http.NewServeMux()
	mux.HandleFunc("/handshake", c.handShakeHandle)
	if c.masterConfig != nil && c.masterConfig.Master {
		mux.HandleFunc("/register", c.registerHandle)
	}

	c.handshakeServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}

	nlog.Error(c.handshakeServer.ListenAndServe())
}

func doHTTP(in interface{}, out interface{}, path, host string, headers map[string]string) error {
	maxRetryTimes := 5

	for i := 0; i < maxRetryTimes; i++ {
		inbody, err := json.Marshal(in)
		if err != nil {
			nlog.Errorf("new handshake request fail:%v", err)
			return err
		}
		req, err := http.NewRequest(http.MethodPost, config.InternalServer+path, bytes.NewBuffer(inbody))
		if err != nil {
			nlog.Errorf("new handshake request fail:%v", err)
			return err
		}
		req.Host = host
		req.Header.Set("Content-Type", "application/json")
		for key, val := range headers {
			req.Header.Set(key, val)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			nlog.Errorf("do http request fail:%v", err)
			time.Sleep(time.Second)
			continue
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			nlog.Warnf("Request error, path: %s, code: %d, message: %s", path, resp.StatusCode, string(body))
			time.Sleep(time.Second)
			continue
		}

		if err = json.Unmarshal(body, out); err != nil {
			nlog.Errorf("Json unmarshal failed, err:%s, body:%s", err.Error(), string(body))
			time.Sleep(time.Second)
			continue
		}
		return nil
	}

	return fmt.Errorf("request error, retry at maxtimes:%d, path: %s", maxRetryTimes, path)
}

func (c *DomainRouteController) waitTokenReady(ctx context.Context, dr *kusciaapisv1alpha1.DomainRoute) error {
	maxRetryTimes := 60
	i := 0
	t := time.NewTicker(time.Second)
	defer t.Stop()
	var err error
	var out *getResponse
	for {
		i++
		if i == maxRetryTimes {
			return fmt.Errorf("wait dr %s token ready failed at max retry times:%d, last error: %s", dr.Name, maxRetryTimes, err.Error())
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			var drLatest *kusciaapisv1alpha1.DomainRoute
			drLatest, err = c.domainRouteLister.DomainRoutes(dr.Namespace).Get(dr.Name)
			if err != nil {
				return err
			}
			if drLatest.Status.TokenStatus.RevisionToken.IsReady {
				return nil
			}
			if drLatest.Status.TokenStatus.RevisionToken.Token == "" {
				return fmt.Errorf("token of dr %s was deleted ", drLatest.Name)
			}
			if drLatest.Status.TokenStatus.RevisionInitializer != c.gateway.Name {
				return fmt.Errorf("dr %s may change initializer ", drLatest.Name)
			}
			out, err = c.checkConnectionStatus(drLatest)
			if err != nil {
				if out != nil && out.State != NetworkUnreachable {
					if out.State == TokenNotReady {
						continue
					} else {
						nlog.Warnf("err:%s, retry time: %d", err.Error(), i)
						return err
					}
				} else {
					nlog.Warnf("err:%s, retry time: %d", err.Error(), i)
				}
			} else {
				nlog.Infof("Destination(%s) token is ready", drLatest.Spec.Destination)
				return nil
			}
		}
	}
}

func (c *DomainRouteController) checkConnectionStatus(dr *kusciaapisv1alpha1.DomainRoute) (*getResponse, error) {
	req, err := http.NewRequest(http.MethodGet, config.InternalServer+"/handshake", nil)
	if err != nil {
		nlog.Errorf("new handshake request fail:%v", err)
		return nil, err
	}
	ns := dr.Spec.Destination
	if dr.Spec.Transit != nil {
		ns = dr.Spec.Transit.Domain.DomainID
	}
	req.Host = fmt.Sprintf("%s.%s.svc", clusters.ServiceHandshake, ns)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(fmt.Sprintf("%s-Cluster", clusters.ServiceHandshake), fmt.Sprintf("%s-to-%s-%s", dr.Spec.Source, ns, dr.Spec.Endpoint.Ports[0].Name))
	req.Header.Set("Kuscia-Token-Revision", fmt.Sprintf("%d", dr.Status.TokenStatus.RevisionToken.Revision))
	req.Header.Set("Kuscia-Source", dr.Spec.Source)
	req.Header.Set("kuscia-Host", fmt.Sprintf("%s.%s.svc", clusters.ServiceHandshake, ns))
	client := &http.Client{}
	resp, err := client.Do(req)

	buildUnreachableResp := func() *getResponse {
		out := &getResponse{
			Namespace: ns,
			State:     NetworkUnreachable,
		}
		return out
	}

	if err != nil {
		return buildUnreachableResp(), fmt.Errorf("do http request fail:%v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return buildUnreachableResp(), fmt.Errorf("read body failed, err:%s, code: %d", err.Error(), resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		if len(body) > 200 {
			body = body[:200]
		}
		return buildUnreachableResp(), fmt.Errorf("request error, path: %s, code: %d, message: %s", fmt.Sprintf("%s.%s.svc/handshake",
			clusters.ServiceHandshake, ns), resp.StatusCode, string(body))
	}
	out := &getResponse{}
	err = json.Unmarshal(body, out)
	if err != nil {
		return buildUnreachableResp(), err
	}

	return out, c.handleGetResponse(out, dr)
}

func (c *DomainRouteController) handleGetResponse(out *getResponse, dr *kusciaapisv1alpha1.DomainRoute) error {
	switch out.State {
	case TokenReady:
		if !dr.Status.TokenStatus.RevisionToken.IsReady {
			dr = dr.DeepCopy()
			dr.Status.TokenStatus.RevisionToken.IsReady = true
			dr.Status.IsDestinationUnreachable = false
			_, err := c.kusciaClient.KusciaV1alpha1().DomainRoutes(dr.Namespace).UpdateStatus(context.Background(), dr, metav1.UpdateOptions{})
			return err
		}
		return nil
	case TokenNotFound:
		if dr.Status.TokenStatus.RevisionToken.Token != "" || dr.Status.TokenStatus.Tokens != nil || dr.Status.TokenStatus.RevisionToken.IsReady {
			dr = dr.DeepCopy()
			dr.Status.TokenStatus = kusciaapisv1alpha1.DomainRouteTokenStatus{}
			_, err := c.kusciaClient.KusciaV1alpha1().DomainRoutes(dr.Namespace).UpdateStatus(context.Background(), dr, metav1.UpdateOptions{})
			return err
		}
		return nil
	case DomainIDInputInvalid:
		return fmt.Errorf("%s destinationreturn  DomainIDInputInvalid, check 'Kuscia-Source' in header", dr.Name)
	case TokenRevisionInputInvalid:
		return fmt.Errorf("%s destination return TokenRevisionInputInvalid, check 'Kuscia-Token-Revision' in header", dr.Name)
	case TokenNotReady:
		return fmt.Errorf("%s destination token is not ready", dr.Name)
	case NoAuthentication:
		if dr.Status.IsDestinationAuthorized {
			dr = dr.DeepCopy()
			dr.Status.IsDestinationAuthorized = false
			dr.Status.IsDestinationUnreachable = false
			dr.Status.TokenStatus = kusciaapisv1alpha1.DomainRouteTokenStatus{}
			c.kusciaClient.KusciaV1alpha1().DomainRoutes(dr.Namespace).UpdateStatus(context.Background(), dr, metav1.UpdateOptions{})
			return fmt.Errorf("%s cant contact destination because destination authentication is false", dr.Name)
		}
		return nil
	case InternalError:
		return fmt.Errorf("%s destination return unkown error", dr.Name)
	default:
		return nil
	}
}

func calcPublicKeyHash(pubStr string) ([]byte, error) {
	srcPub, err := base64.StdEncoding.DecodeString(pubStr)
	if err != nil {
		return nil, err
	}
	msgHash := sha256.New()
	_, err = msgHash.Write(srcPub)
	if err != nil {
		return nil, err
	}
	return msgHash.Sum(nil), nil
}

func (c *DomainRouteController) sourceInitiateHandShake(dr *kusciaapisv1alpha1.DomainRoute) error {
	if dr.Spec.TokenConfig.SourcePublicKey != c.gateway.Status.PublicKey {
		nlog.Errorf("DomainRoute %s: mismatch source public key", dr.Name)
		return nil
	}

	handshankeReq := &handshake.HandShakeRequest{
		DomainId:    dr.Spec.Source,
		RequestTime: time.Now().UnixNano(),
	}

	//1. In UID mode, the token is directly generated by the peer end and encrypted by the local public key
	//2. In RSA mode, the local end and the peer end generate their own tokens and concatenate them.
	//   The local token is encrypted with the peer's public key and then sent.
	//   The peer token is encrypted with the local public key and returned.
	var token []byte
	resp := &handshake.HandShakeResponse{}
	if dr.Spec.TokenConfig.TokenGenMethod == kusciaapisv1alpha1.TokenGenUIDRSA {
		handshankeReq.Type = handShakeTypeUID
		ns := dr.Spec.Destination
		if dr.Spec.Transit != nil {
			ns = dr.Spec.Transit.Domain.DomainID
		}
		headers := map[string]string{
			fmt.Sprintf("%s-Cluster", clusters.ServiceHandshake): fmt.Sprintf("%s-to-%s-%s", dr.Spec.Source, ns, dr.Spec.Endpoint.Ports[0].Name),
			"Kuscia-Source": dr.Spec.Source,
			"kuscia-Host":   fmt.Sprintf("%s.%s.svc", clusters.ServiceHandshake, ns),
		}
		err := doHTTP(handshankeReq, resp, "/handshake", fmt.Sprintf("%s.%s.svc", clusters.ServiceHandshake, ns), headers)
		if err != nil {
			nlog.Errorf("DomainRoute %s: handshake fail:%v", dr.Name, err)
			return err
		}
		if resp.Status.Code != 0 {
			err = fmt.Errorf("DomainRoute %s: handshake fail, return error:%v", dr.Name, resp.Status.Message)
			nlog.Error(err)
			return err
		}
		token, err = decryptToken(c.prikey, resp.Token.Token, tokenByteSize)
		if err != nil {
			err = fmt.Errorf("DomainRoute %s: handshake fail, return error:%v", dr.Name, resp.Status.Message)
			nlog.Error(err)
			return err
		}
	} else if dr.Spec.TokenConfig.TokenGenMethod == kusciaapisv1alpha1.TokenGenMethodRSA {
		handshankeReq.Type = handShakeTypeRSA

		msgHashSum, err := calcPublicKeyHash(c.gateway.Status.PublicKey)
		if err != nil {
			return err
		}

		sourceToken := make([]byte, tokenByteSize/2)
		_, err = rand.Read(sourceToken)
		if err != nil {
			return err
		}

		//Resolve the public key of the peer end from domainroute crd
		destPub, err := base64.StdEncoding.DecodeString(dr.Spec.TokenConfig.DestinationPublicKey)
		if err != nil {
			nlog.Errorf("DomainRoute %s: destination public key format error, must be base64 encoded", dr.Name)
			return err
		}
		destPubKey, err := tlsutils.ParsePKCS1PublicKey(destPub)
		if err != nil {
			return err
		}

		sourceTokenEnc, err := encryptToken(destPubKey, sourceToken)
		if err != nil {
			return err
		}

		handshankeReq.TokenConfig = &handshake.TokenConfig{
			Token:    sourceTokenEnc,
			Revision: dr.Status.TokenStatus.RevisionToken.Revision,
			Pubhash:  base64.StdEncoding.EncodeToString(msgHashSum),
		}

		ns := dr.Spec.Destination
		if dr.Spec.Transit != nil {
			ns = dr.Spec.Transit.Domain.DomainID
		}
		headers := map[string]string{
			fmt.Sprintf("%s-Cluster", clusters.ServiceHandshake): fmt.Sprintf("%s-to-%s-%s", dr.Spec.Source, ns, dr.Spec.Endpoint.Ports[0].Name),
			"Kuscia-Source": dr.Spec.Source,
			"kuscia-Host":   fmt.Sprintf("%s.%s.svc", clusters.ServiceHandshake, ns),
		}

		err = doHTTP(handshankeReq, resp, "/handshake", fmt.Sprintf("%s.%s.svc", clusters.ServiceHandshake, ns), headers)
		if err != nil {
			nlog.Errorf("DomainRoute %s: handshake fail:%v", dr.Name, err)
			return err
		}
		if resp.Status.Code != 0 {
			err = fmt.Errorf("DomainRoute %s: handshake fail, return error:%v", dr.Name, resp.Status.Message)
			nlog.Error(err)
			return err
		}
		destToken, err := decryptToken(c.prikey, resp.Token.Token, tokenByteSize/2)
		if err != nil {
			err = fmt.Errorf("DomainRoute %s: handshake fail, decryptToken  error:%v", dr.Name, resp.Status.Message)
			nlog.Error(err)
			return err
		}
		token = append(sourceToken, destToken...)
	} else {
		return fmt.Errorf("TokenGenMethod must be %s or %s", kusciaapisv1alpha1.TokenGenUIDRSA, kusciaapisv1alpha1.TokenGenMethodRSA)
	}

	// The final token is encrypted with the local private key and stored in the status of domainroute
	tokenEncrypted, err := encryptToken(&c.prikey.PublicKey, token)
	if err != nil {
		return err
	}
	drLatest, _ := c.domainRouteLister.DomainRoutes(dr.Namespace).Get(dr.Name)
	drCopy := drLatest.DeepCopy()
	tn := metav1.Now()
	drCopy.Status.IsDestinationAuthorized = true
	drCopy.Status.TokenStatus.RevisionToken.Token = tokenEncrypted
	drCopy.Status.TokenStatus.RevisionToken.Revision = int64(resp.Token.Revision)
	drCopy.Status.TokenStatus.RevisionToken.IsReady = false
	drCopy.Status.TokenStatus.RevisionToken.RevisionTime = tn
	if drCopy.Spec.TokenConfig.RollingUpdatePeriod == 0 {
		drCopy.Status.TokenStatus.RevisionToken.ExpirationTime = metav1.NewTime(tn.AddDate(100, 0, 0))
	} else {
		tTx := time.Unix(resp.Token.ExpirationTime/int64(time.Second), resp.Token.ExpirationTime%int64(time.Second))
		drCopy.Status.TokenStatus.RevisionToken.ExpirationTime = metav1.NewTime(tTx)
	}
	_, err = c.kusciaClient.KusciaV1alpha1().DomainRoutes(drCopy.Namespace).UpdateStatus(context.Background(), drCopy, metav1.UpdateOptions{})
	return err
}

type DestinationStatus int

const (
	TokenReady DestinationStatus = iota
	DomainIDInputInvalid
	TokenRevisionInputInvalid
	TokenNotReady
	TokenNotFound
	NoAuthentication
	InternalError
	NetworkUnreachable
)

type getResponse struct {
	Namespace string            `json:"namespace"`
	State     DestinationStatus `json:"state"`
}

func (c *DomainRouteController) handShakeHandle(w http.ResponseWriter, r *http.Request) {
	nlog.Debugf("Receive handshake request, method [%s], host[%s], headers[%s]", r.Method, r.Host, r.Header)
	if r.Method == http.MethodGet {
		resp := &getResponse{
			Namespace: c.gateway.Namespace,
			State:     TokenNotReady,
		}

		domainID := r.Header.Get("Kuscia-Source")
		tokenRevision := r.Header.Get("Kuscia-Token-Revision")
		resp.State = c.checkTokenStatus(domainID, tokenRevision)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			nlog.Errorf("write handshake response fail, detail-> %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	req := handshake.HandShakeRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		nlog.Errorf("Invalid request: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	drName := common.GenDomainRouteName(req.DomainId, c.gateway.Namespace)
	dr, err := c.domainRouteLister.DomainRoutes(c.gateway.Namespace).Get(drName)
	if err != nil {
		msg := fmt.Sprintf("DomainRoute %s get error: %v", drName, err)
		nlog.Error(msg)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	if !(req.Type == handShakeTypeUID && dr.Spec.TokenConfig.TokenGenMethod == kusciaapisv1alpha1.TokenGenUIDRSA) &&
		!(req.Type == handShakeTypeRSA && dr.Spec.TokenConfig.TokenGenMethod == kusciaapisv1alpha1.TokenGenMethodRSA) {
		errMsg := fmt.Sprintf("handshake Type(%s) not match domainroute required(%s)", req.Type, dr.Spec.TokenConfig.TokenGenMethod)
		nlog.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	resp := c.DestReplyHandshake(&req, dr)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		nlog.Errorf("encode handshake response for(%s) fail, detail-> %v", drName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if resp.Status.Code != 0 {
			nlog.Errorf("DestReplyHandshake for(%s) fail, detail-> %v", drName, resp.Status.Message)
		} else {
			nlog.Infof("DomainRoute %s handle success", drName)
		}
	}
}

func (c *DomainRouteController) checkTokenStatus(domainID, tokenRevision string) DestinationStatus {
	if domainID == "" {
		return DomainIDInputInvalid
	}
	drName := common.GenDomainRouteName(domainID, c.gateway.Namespace)
	if dr, err := c.domainRouteLister.DomainRoutes(c.gateway.Namespace).Get(drName); err == nil {
		return checkTokenRevision(tokenRevision, dr)
	} else if k8serrors.IsNotFound(err) {
		return NoAuthentication
	}
	return InternalError
}

func checkTokenRevision(tokenRevision string, dr *kusciaapisv1alpha1.DomainRoute) DestinationStatus {
	if tokenRevision == "" {
		return TokenRevisionInputInvalid
	}
	r, err := strconv.Atoi(tokenRevision)
	if err != nil {
		return TokenRevisionInputInvalid
	}
	for _, t := range dr.Status.TokenStatus.Tokens {
		if t.Revision == int64(r) {
			if t.IsReady {
				return TokenReady
			}
			return TokenNotReady
		}
	}
	return TokenNotFound
}

func buildFailedHandshakeReply(code int32, err error) *handshake.HandShakeResponse {
	resp := &handshake.HandShakeResponse{
		Status: &v1alpha1.Status{
			Code: code,
		},
	}
	if err != nil {
		resp.Status.Message = err.Error()
	}
	return resp
}

func (c *DomainRouteController) DestReplyHandshake(req *handshake.HandShakeRequest, dr *kusciaapisv1alpha1.DomainRoute) *handshake.HandShakeResponse {
	srcPub, err := base64.StdEncoding.DecodeString(dr.Spec.TokenConfig.SourcePublicKey)
	if err != nil {
		return buildFailedHandshakeReply(500, err)
	}
	sourcePubKey, err := tlsutils.ParsePKCS1PublicKey(srcPub)
	if err != nil {
		return buildFailedHandshakeReply(500, err)
	}
	dstRevisionToken := dr.Status.TokenStatus.RevisionToken

	var token []byte
	var respToken []byte
	if req.Type == handShakeTypeUID {
		// If the token in domainroute is empty or has expired, the token is regenerated.
		// Otherwise, the token is returned
		needGenerateToken := func() bool {
			if dstRevisionToken.Token == "" {
				return true
			}
			if dr.Spec.TokenConfig != nil && dr.Spec.TokenConfig.RollingUpdatePeriod > 0 && time.Since(dstRevisionToken.RevisionTime.Time) > time.Duration(dr.Spec.TokenConfig.RollingUpdatePeriod)*time.Second {
				return true
			}
			return false
		}
		if needGenerateToken() {
			respToken, err = generateRandomToken(tokenByteSize)
			if err != nil {
				return buildFailedHandshakeReply(500, err)
			}
		} else {
			respToken, err = decryptToken(c.prikey, dstRevisionToken.Token, tokenByteSize)
			if err != nil {
				nlog.Warnf("source %s %s handshake decryptToken failed, error:%s", req.DomainId, handShakeTypeUID, err.Error())
				respToken, err = generateRandomToken(tokenByteSize)
				if err != nil {
					return buildFailedHandshakeReply(500, err)
				}
			}
		}

		token = respToken
	} else if req.Type == handShakeTypeRSA {
		msgHashSum, err := calcPublicKeyHash(dr.Spec.TokenConfig.SourcePublicKey)
		if err != nil {
			return buildFailedHandshakeReply(500, fmt.Errorf("source %s %s handshake calcPublicKeyHash failed, error:%s", req.DomainId, req.Type, err.Error()))
		}
		if req.TokenConfig.Pubhash != base64.StdEncoding.EncodeToString(msgHashSum) {
			return buildFailedHandshakeReply(500, fmt.Errorf("source %s %s publickey mismatch in domainroute(%s) SourcePublicKey", req.DomainId, req.Type, dr.Name))
		}
		sourceToken, err := decryptToken(c.prikey, req.TokenConfig.Token, tokenByteSize/2)
		if err != nil {
			return buildFailedHandshakeReply(500, fmt.Errorf("source %s %s handshake decryptToken failed, error:%s", req.DomainId, req.Type, err.Error()))
		}

		respToken = make([]byte, tokenByteSize/2)
		if _, err = rand.Read(respToken); err != nil {
			return buildFailedHandshakeReply(500, err)
		}

		token = append(sourceToken, respToken...)
	}

	tokenEncrypted, err := encryptToken(&c.prikey.PublicKey, token)
	if err != nil {
		return buildFailedHandshakeReply(500, fmt.Errorf("source %s %s handshake encryptToken by self publickey failed, error:%s", req.DomainId, req.Type, err.Error()))
	}

	respTokenEncrypted, err := encryptToken(sourcePubKey, respToken)
	if err != nil {
		return buildFailedHandshakeReply(500, fmt.Errorf("source %s %s handshake encryptToken by source publickey failed, error:%s", req.DomainId, req.Type, err.Error()))
	}

	drLatest, err := c.domainRouteLister.DomainRoutes(dr.Namespace).Get(dr.Name)
	if err != nil {
		return buildFailedHandshakeReply(500, err)
	}
	var revision int64
	var expirationTime metav1.Time
	if drLatest.Status.TokenStatus.RevisionToken.Token != tokenEncrypted {
		drCopy := drLatest.DeepCopy()
		revisionTime := metav1.Now()
		drCopy.Status.TokenStatus.RevisionToken.Token = tokenEncrypted
		drCopy.Status.TokenStatus.RevisionToken.Revision++
		revision = drCopy.Status.TokenStatus.RevisionToken.Revision
		drCopy.Status.TokenStatus.RevisionToken.RevisionTime = revisionTime
		drCopy.Status.TokenStatus.RevisionToken.IsReady = false
		if drCopy.Spec.TokenConfig.RollingUpdatePeriod == 0 {
			expirationTime = metav1.NewTime(revisionTime.AddDate(100, 0, 0))
		} else {
			expirationTime = metav1.NewTime(revisionTime.Add(2 * time.Duration(drCopy.Spec.TokenConfig.RollingUpdatePeriod) * time.Second))
		}
		drCopy.Status.TokenStatus.RevisionToken.ExpirationTime = expirationTime
		_, err = c.kusciaClient.KusciaV1alpha1().DomainRoutes(drCopy.Namespace).UpdateStatus(context.Background(), drCopy, metav1.UpdateOptions{})
		if err != nil {
			return buildFailedHandshakeReply(500, err)
		}
		nlog.Infof("Update domainroute %s status", dr.Name)
	} else {
		revision = drLatest.Status.TokenStatus.RevisionToken.Revision
		expirationTime = drLatest.Status.TokenStatus.RevisionToken.ExpirationTime
	}

	return &handshake.HandShakeResponse{
		Status: &v1alpha1.Status{
			Code: 0,
		},
		Token: &handshake.Token{
			Token:          respTokenEncrypted,
			ExpirationTime: expirationTime.UnixNano(),
			Revision:       int32(revision),
		},
	}
}

func (c *DomainRouteController) parseToken(dr *kusciaapisv1alpha1.DomainRoute, routeKey string) ([]*Token, error) {
	var tokens []*Token
	var err error

	if (dr.Spec.Transit != nil && dr.Spec.BodyEncryption == nil) ||
		(dr.Spec.Transit == nil && dr.Spec.AuthenticationType == kusciaapisv1alpha1.DomainAuthenticationMTLS) ||
		(dr.Spec.Transit == nil && dr.Spec.AuthenticationType == kusciaapisv1alpha1.DomainAuthenticationNone) {
		tokens = append(tokens, &Token{Token: NoopToken})
		return tokens, err
	}

	if (dr.Spec.Transit == nil && dr.Spec.AuthenticationType != kusciaapisv1alpha1.DomainAuthenticationToken) ||
		dr.Spec.TokenConfig == nil {
		return tokens, fmt.Errorf("invalid DomainRoute: %v", dr.Spec)
	}

	switch dr.Spec.TokenConfig.TokenGenMethod {
	case kusciaapisv1alpha1.TokenGenMethodRSA:
		tokens, err = c.parseTokenRSA(dr, false)
	case kusciaapisv1alpha1.TokenGenUIDRSA:
		tokens, err = c.parseTokenRSA(dr, true)
	default:
		err = fmt.Errorf("DomainRoute %s unsupported token method: %s", routeKey,
			dr.Spec.TokenConfig.TokenGenMethod)
	}
	return tokens, err
}

func (c *DomainRouteController) parseTokenRSA(dr *kusciaapisv1alpha1.DomainRoute, drop bool) ([]*Token, error) {
	key, _ := cache.MetaNamespaceKeyFunc(dr)

	var tokens []*Token
	if len(dr.Status.TokenStatus.Tokens) == 0 {
		return tokens, nil
	}

	if (c.gateway.Namespace == dr.Spec.Source && dr.Spec.TokenConfig.SourcePublicKey != c.gateway.Status.PublicKey) ||
		(c.gateway.Namespace == dr.Spec.Destination && dr.Spec.TokenConfig.DestinationPublicKey != c.gateway.Status.PublicKey) {
		err := fmt.Errorf("DomainRoute %s mismatch public key", key)
		return tokens, err
	}

	for _, token := range dr.Status.TokenStatus.Tokens {
		b, err := decryptToken(c.prikey, token.Token, tokenByteSize)
		if err != nil {
			if !drop {
				return []*Token{}, fmt.Errorf("DomainRoute %s decrypt token error: %v", key, err)
			}
			nlog.Warnf("DomainRoute %s decrypt token [revision -> %d] error: %v", key, token.Revision, err)
			continue
		}
		tokens = append(tokens, &Token{Token: base64.StdEncoding.EncodeToString(b), Version: token.Revision})
	}

	return tokens, nil
}

func (c *DomainRouteController) checkAndUpdateTokenInstances(dr *kusciaapisv1alpha1.DomainRoute) error {
	if len(dr.Status.TokenStatus.Tokens) == 0 {
		return nil
	}
	updated := false
	drCopy := dr.DeepCopy()
	for i := range drCopy.Status.TokenStatus.Tokens {
		if !exists(drCopy.Status.TokenStatus.Tokens[i].EffectiveInstances, c.gateway.Name) {
			updated = true
			drCopy.Status.TokenStatus.Tokens[i].EffectiveInstances = append(drCopy.Status.TokenStatus.Tokens[i].EffectiveInstances, c.gateway.Name)
		}
	}
	if updated {
		_, err := c.kusciaClient.KusciaV1alpha1().DomainRoutes(drCopy.Namespace).UpdateStatus(context.Background(), drCopy, metav1.UpdateOptions{})
		return err
	}
	return nil
}

func generateRandomToken(size int) ([]byte, error) {
	respToken := make([]byte, size)
	if _, err := rand.Read(respToken); err != nil {
		return nil, err
	}
	return respToken, nil
}

func encryptToken(pub *rsa.PublicKey, key []byte) (string, error) {
	return tlsutils.EncryptPKCS1v15(pub, key, tokenPrefix)
}

func decryptToken(priv *rsa.PrivateKey, ciphertext string, keysize int) ([]byte, error) {
	return tlsutils.DecryptPKCS1v15(priv, ciphertext, keysize, tokenPrefix)
}

func exists(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func HandshakeToMaster(domainID string, prikey *rsa.PrivateKey) error {
	handshankeReq := &handshake.HandShakeRequest{
		DomainId:    domainID,
		RequestTime: time.Now().UnixNano(),
	}

	//1. In UID mode, the token is directly generated by the peer end and encrypted by the local public key
	//2. In RSA mode, the local end and the peer end generate their own tokens and concatenate them.
	//   The local token is encrypted with the peer's public key and then sent.
	//   The peer token is encrypted with the local public key and returned.
	handshankeReq.Type = handShakeTypeUID
	resp := &handshake.HandShakeResponse{}

	headers := map[string]string{
		"Kuscia-Source": domainID,
		fmt.Sprintf("%s-Cluster", clusters.ServiceHandshake): clusters.GetMasterClusterName(),
		"kuscia-Host": fmt.Sprintf("%s.master.svc", clusters.ServiceHandshake),
	}
	err := doHTTP(handshankeReq, resp, "/handshake", fmt.Sprintf("%s.master.svc", clusters.ServiceHandshake), headers)
	if err != nil {
		nlog.Error(err)
		return err
	}
	if resp.Status.Code != 0 {
		nlog.Errorf("Handshake to master fail, return error:%v", resp.Status.Message)
		return errors.New(resp.Status.Message)
	}
	token, err := decryptToken(prikey, resp.Token.Token, tokenByteSize)
	if err != nil {
		nlog.Errorf("Handshake to master decryptToken err:%s", err.Error())
		return err
	}
	c, err := xds.QueryCluster(clusters.GetMasterClusterName())
	if err != nil {
		nlog.Error(err)
		return err
	}
	if err := clusters.AddMasterProxyVirtualHost(c.Name, clusters.ServiceMasterProxy, domainID, base64.StdEncoding.EncodeToString(token)); err != nil {
		nlog.Error(err)
		return err
	}
	if err = xds.SetKeepAliveForDstCluster(c, true); err != nil {
		nlog.Error(err)
		return err
	}
	if err = xds.AddOrUpdateCluster(c); err != nil {
		nlog.Error(err)
		return err
	}
	return nil
}
