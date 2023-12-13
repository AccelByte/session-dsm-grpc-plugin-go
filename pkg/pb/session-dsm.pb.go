// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.18.1
// source: session-dsm.proto

package sessiondsm

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RequestCreateGameSession struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId       string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Namespace       string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	FleetAlias      string `protobuf:"bytes,3,opt,name=fleet_alias,json=fleetAlias,proto3" json:"fleet_alias,omitempty"`
	SessionData     string `protobuf:"bytes,4,opt,name=session_data,json=sessionData,proto3" json:"session_data,omitempty"`
	RequestedRegion string `protobuf:"bytes,5,opt,name=requested_region,json=requestedRegion,proto3" json:"requested_region,omitempty"`
	MaximumPlayer   int64  `protobuf:"varint,6,opt,name=maximum_player,json=maximumPlayer,proto3" json:"maximum_player,omitempty"`
}

func (x *RequestCreateGameSession) Reset() {
	*x = RequestCreateGameSession{}
	if protoimpl.UnsafeEnabled {
		mi := &file_session_dsm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestCreateGameSession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestCreateGameSession) ProtoMessage() {}

func (x *RequestCreateGameSession) ProtoReflect() protoreflect.Message {
	mi := &file_session_dsm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestCreateGameSession.ProtoReflect.Descriptor instead.
func (*RequestCreateGameSession) Descriptor() ([]byte, []int) {
	return file_session_dsm_proto_rawDescGZIP(), []int{0}
}

func (x *RequestCreateGameSession) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *RequestCreateGameSession) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *RequestCreateGameSession) GetFleetAlias() string {
	if x != nil {
		return x.FleetAlias
	}
	return ""
}

func (x *RequestCreateGameSession) GetSessionData() string {
	if x != nil {
		return x.SessionData
	}
	return ""
}

func (x *RequestCreateGameSession) GetRequestedRegion() string {
	if x != nil {
		return x.RequestedRegion
	}
	return ""
}

func (x *RequestCreateGameSession) GetMaximumPlayer() int64 {
	if x != nil {
		return x.MaximumPlayer
	}
	return 0
}

type ResponseCreateGameSession struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId     string                 `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Namespace     string                 `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	FleetAlias    string                 `protobuf:"bytes,3,opt,name=fleet_alias,json=fleetAlias,proto3" json:"fleet_alias,omitempty"`
	SessionData   string                 `protobuf:"bytes,4,opt,name=session_data,json=sessionData,proto3" json:"session_data,omitempty"`
	Status        string                 `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	Ip            string                 `protobuf:"bytes,6,opt,name=ip,proto3" json:"ip,omitempty"`
	Port          int64                  `protobuf:"varint,7,opt,name=port,proto3" json:"port,omitempty"`
	ServerId      string                 `protobuf:"bytes,8,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`
	Source        string                 `protobuf:"bytes,9,opt,name=source,proto3" json:"source,omitempty"`
	Deployment    string                 `protobuf:"bytes,10,opt,name=deployment,proto3" json:"deployment,omitempty"`
	Region        string                 `protobuf:"bytes,11,opt,name=region,proto3" json:"region,omitempty"`
	LastUpdatedAt *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=last_updated_at,json=lastUpdatedAt,proto3" json:"last_updated_at,omitempty"`
}

func (x *ResponseCreateGameSession) Reset() {
	*x = ResponseCreateGameSession{}
	if protoimpl.UnsafeEnabled {
		mi := &file_session_dsm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseCreateGameSession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseCreateGameSession) ProtoMessage() {}

func (x *ResponseCreateGameSession) ProtoReflect() protoreflect.Message {
	mi := &file_session_dsm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseCreateGameSession.ProtoReflect.Descriptor instead.
func (*ResponseCreateGameSession) Descriptor() ([]byte, []int) {
	return file_session_dsm_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseCreateGameSession) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *ResponseCreateGameSession) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ResponseCreateGameSession) GetFleetAlias() string {
	if x != nil {
		return x.FleetAlias
	}
	return ""
}

func (x *ResponseCreateGameSession) GetSessionData() string {
	if x != nil {
		return x.SessionData
	}
	return ""
}

func (x *ResponseCreateGameSession) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ResponseCreateGameSession) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *ResponseCreateGameSession) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *ResponseCreateGameSession) GetServerId() string {
	if x != nil {
		return x.ServerId
	}
	return ""
}

func (x *ResponseCreateGameSession) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *ResponseCreateGameSession) GetDeployment() string {
	if x != nil {
		return x.Deployment
	}
	return ""
}

func (x *ResponseCreateGameSession) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *ResponseCreateGameSession) GetLastUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUpdatedAt
	}
	return nil
}

var File_session_dsm_proto protoreflect.FileDescriptor

var file_session_dsm_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x64, 0x73, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x61, 0x63, 0x63, 0x65, 0x6c, 0x62, 0x79, 0x74, 0x65, 0x2e, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x64, 0x73,
	0x6d, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xed, 0x01, 0x0a, 0x18, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x66, 0x6c, 0x65, 0x65, 0x74, 0x5f, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x21, 0x0a,
	0x0c, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x29, 0x0a, 0x10, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x65,
	0x67, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x6d,
	0x61, 0x78, 0x69, 0x6d, 0x75, 0x6d, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0d, 0x6d, 0x61, 0x78, 0x69, 0x6d, 0x75, 0x6d, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x22, 0x89, 0x03, 0x0a, 0x19, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x5f, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x21,
	0x0a, 0x0c, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x42, 0x0a, 0x0f, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0d, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x32, 0x93,
	0x01, 0x0a, 0x0a, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x73, 0x6d, 0x12, 0x84, 0x01,
	0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x36, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x6c, 0x62, 0x79, 0x74, 0x65, 0x2e,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x64,
	0x73, 0x6d, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x47, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x37, 0x2e, 0x61, 0x63,
	0x63, 0x65, 0x6c, 0x62, 0x79, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x64, 0x73, 0x6d, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x42, 0x65, 0x0a, 0x20, 0x6e, 0x65, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x65,
	0x6c, 0x62, 0x79, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x64, 0x73, 0x6d, 0x50, 0x01, 0x5a, 0x20, 0x61, 0x63, 0x63, 0x65,
	0x6c, 0x62, 0x79, 0x74, 0x65, 0x2e, 0x6e, 0x65, 0x74, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x64, 0x73, 0x6d, 0xaa, 0x02, 0x1c, 0x41,
	0x63, 0x63, 0x65, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x64, 0x73, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_session_dsm_proto_rawDescOnce sync.Once
	file_session_dsm_proto_rawDescData = file_session_dsm_proto_rawDesc
)

func file_session_dsm_proto_rawDescGZIP() []byte {
	file_session_dsm_proto_rawDescOnce.Do(func() {
		file_session_dsm_proto_rawDescData = protoimpl.X.CompressGZIP(file_session_dsm_proto_rawDescData)
	})
	return file_session_dsm_proto_rawDescData
}

var file_session_dsm_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_session_dsm_proto_goTypes = []interface{}{
	(*RequestCreateGameSession)(nil),  // 0: accelbyte.session.sessiondsm.RequestCreateGameSession
	(*ResponseCreateGameSession)(nil), // 1: accelbyte.session.sessiondsm.ResponseCreateGameSession
	(*timestamppb.Timestamp)(nil),     // 2: google.protobuf.Timestamp
}
var file_session_dsm_proto_depIdxs = []int32{
	2, // 0: accelbyte.session.sessiondsm.ResponseCreateGameSession.last_updated_at:type_name -> google.protobuf.Timestamp
	0, // 1: accelbyte.session.sessiondsm.SessionDsm.CreateGameSession:input_type -> accelbyte.session.sessiondsm.RequestCreateGameSession
	1, // 2: accelbyte.session.sessiondsm.SessionDsm.CreateGameSession:output_type -> accelbyte.session.sessiondsm.ResponseCreateGameSession
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_session_dsm_proto_init() }
func file_session_dsm_proto_init() {
	if File_session_dsm_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_session_dsm_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestCreateGameSession); i {
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
		file_session_dsm_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseCreateGameSession); i {
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
			RawDescriptor: file_session_dsm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_session_dsm_proto_goTypes,
		DependencyIndexes: file_session_dsm_proto_depIdxs,
		MessageInfos:      file_session_dsm_proto_msgTypes,
	}.Build()
	File_session_dsm_proto = out.File
	file_session_dsm_proto_rawDesc = nil
	file_session_dsm_proto_goTypes = nil
	file_session_dsm_proto_depIdxs = nil
}
