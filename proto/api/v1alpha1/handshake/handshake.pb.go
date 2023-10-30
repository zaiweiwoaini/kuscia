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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: kuscia/proto/api/v1alpha1/handshake/handshake.proto

package handshake

import (
	v1alpha1 "github.com/secretflow/kuscia/proto/api/v1alpha1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TokenConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token    string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Revision int64  `protobuf:"varint,2,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (x *TokenConfig) Reset() {
	*x = TokenConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenConfig) ProtoMessage() {}

func (x *TokenConfig) ProtoReflect() protoreflect.Message {
	mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenConfig.ProtoReflect.Descriptor instead.
func (*TokenConfig) Descriptor() ([]byte, []int) {
	return file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescGZIP(), []int{0}
}

func (x *TokenConfig) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *TokenConfig) GetRevision() int64 {
	if x != nil {
		return x.Revision
	}
	return 0
}

// HandShakeRequest represents domainroute handshake request.
type HandShakeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DomainID.
	DomainId string `protobuf:"bytes,1,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	// Type of handshake.
	Type        string       `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	TokenConfig *TokenConfig `protobuf:"bytes,3,opt,name=token_config,json=tokenConfig,proto3" json:"token_config,omitempty"`
	RequestTime int64        `protobuf:"varint,4,opt,name=request_time,json=requestTime,proto3" json:"request_time,omitempty"`
}

func (x *HandShakeRequest) Reset() {
	*x = HandShakeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandShakeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandShakeRequest) ProtoMessage() {}

func (x *HandShakeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandShakeRequest.ProtoReflect.Descriptor instead.
func (*HandShakeRequest) Descriptor() ([]byte, []int) {
	return file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescGZIP(), []int{1}
}

func (x *HandShakeRequest) GetDomainId() string {
	if x != nil {
		return x.DomainId
	}
	return ""
}

func (x *HandShakeRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *HandShakeRequest) GetTokenConfig() *TokenConfig {
	if x != nil {
		return x.TokenConfig
	}
	return nil
}

func (x *HandShakeRequest) GetRequestTime() int64 {
	if x != nil {
		return x.RequestTime
	}
	return 0
}

type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token          string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ExpirationTime int64  `protobuf:"varint,2,opt,name=expiration_time,json=expirationTime,proto3" json:"expiration_time,omitempty"`
	Revision       int32  `protobuf:"varint,3,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescGZIP(), []int{2}
}

func (x *Token) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Token) GetExpirationTime() int64 {
	if x != nil {
		return x.ExpirationTime
	}
	return 0
}

func (x *Token) GetRevision() int32 {
	if x != nil {
		return x.Revision
	}
	return 0
}

type HandShakeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *v1alpha1.Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Token  *Token           `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *HandShakeResponse) Reset() {
	*x = HandShakeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandShakeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandShakeResponse) ProtoMessage() {}

func (x *HandShakeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandShakeResponse.ProtoReflect.Descriptor instead.
func (*HandShakeResponse) Descriptor() ([]byte, []int) {
	return file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescGZIP(), []int{3}
}

func (x *HandShakeResponse) GetStatus() *v1alpha1.Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *HandShakeResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DomainId    string `protobuf:"bytes,1,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	Csr         string `protobuf:"bytes,2,opt,name=csr,proto3" json:"csr,omitempty"`
	RequestTime int64  `protobuf:"varint,3,opt,name=request_time,json=requestTime,proto3" json:"request_time,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterRequest) GetDomainId() string {
	if x != nil {
		return x.DomainId
	}
	return ""
}

func (x *RegisterRequest) GetCsr() string {
	if x != nil {
		return x.Csr
	}
	return ""
}

func (x *RegisterRequest) GetRequestTime() int64 {
	if x != nil {
		return x.RequestTime
	}
	return 0
}

type RegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *v1alpha1.Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Cert   string           `protobuf:"bytes,3,opt,name=cert,proto3" json:"cert,omitempty"`
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescGZIP(), []int{5}
}

func (x *RegisterResponse) GetStatus() *v1alpha1.Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *RegisterResponse) GetCert() string {
	if x != nil {
		return x.Cert
	}
	return ""
}

var File_kuscia_proto_api_v1alpha1_handshake_handshake_proto protoreflect.FileDescriptor

var file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDesc = []byte{
	0x0a, 0x33, 0x6b, 0x75, 0x73, 0x63, 0x69, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x68, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x23, 0x6b, 0x75, 0x73, 0x63, 0x69, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x1a, 0x26, 0x6b, 0x75, 0x73, 0x63,
	0x69, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x0b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0xbb, 0x01, 0x0a, 0x10, 0x48, 0x61, 0x6e, 0x64, 0x53, 0x68, 0x61, 0x6b,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x53, 0x0a, 0x0c, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x30, 0x2e, 0x6b, 0x75, 0x73, 0x63, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x68, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x21,
	0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x22, 0x62, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x72, 0x65, 0x76,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x90, 0x01, 0x0a, 0x11, 0x48, 0x61, 0x6e, 0x64, 0x53, 0x68,
	0x61, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6b, 0x75,
	0x73, 0x63, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x40, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6b, 0x75, 0x73, 0x63, 0x69, 0x61, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2e, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x63, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x73, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x73, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x61, 0x0a,
	0x10, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x39, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x21, 0x2e, 0x6b, 0x75, 0x73, 0x63, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x65, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74,
	0x42, 0x5e, 0x0a, 0x21, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x66, 0x6c,
	0x6f, 0x77, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x68, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x6b, 0x75, 0x73,
	0x63, 0x69, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescOnce sync.Once
	file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescData = file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDesc
)

func file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescGZIP() []byte {
	file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescOnce.Do(func() {
		file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescData = protoimpl.X.CompressGZIP(file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescData)
	})
	return file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDescData
}

var file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_goTypes = []interface{}{
	(*TokenConfig)(nil),       // 0: kuscia.proto.api.v1alpha1.handshake.TokenConfig
	(*HandShakeRequest)(nil),  // 1: kuscia.proto.api.v1alpha1.handshake.HandShakeRequest
	(*Token)(nil),             // 2: kuscia.proto.api.v1alpha1.handshake.Token
	(*HandShakeResponse)(nil), // 3: kuscia.proto.api.v1alpha1.handshake.HandShakeResponse
	(*RegisterRequest)(nil),   // 4: kuscia.proto.api.v1alpha1.handshake.RegisterRequest
	(*RegisterResponse)(nil),  // 5: kuscia.proto.api.v1alpha1.handshake.RegisterResponse
	(*v1alpha1.Status)(nil),   // 6: kuscia.proto.api.v1alpha1.Status
}
var file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_depIdxs = []int32{
	0, // 0: kuscia.proto.api.v1alpha1.handshake.HandShakeRequest.token_config:type_name -> kuscia.proto.api.v1alpha1.handshake.TokenConfig
	6, // 1: kuscia.proto.api.v1alpha1.handshake.HandShakeResponse.status:type_name -> kuscia.proto.api.v1alpha1.Status
	2, // 2: kuscia.proto.api.v1alpha1.handshake.HandShakeResponse.token:type_name -> kuscia.proto.api.v1alpha1.handshake.Token
	6, // 3: kuscia.proto.api.v1alpha1.handshake.RegisterResponse.status:type_name -> kuscia.proto.api.v1alpha1.Status
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_init() }
func file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_init() {
	if File_kuscia_proto_api_v1alpha1_handshake_handshake_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandShakeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Token); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandShakeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_goTypes,
		DependencyIndexes: file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_depIdxs,
		MessageInfos:      file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_msgTypes,
	}.Build()
	File_kuscia_proto_api_v1alpha1_handshake_handshake_proto = out.File
	file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_rawDesc = nil
	file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_goTypes = nil
	file_kuscia_proto_api_v1alpha1_handshake_handshake_proto_depIdxs = nil
}