// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proxy/proxy.proto

package proxy

import (
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

type RegisterNodeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RateLimiterId int64 `protobuf:"varint,1,opt,name=rate_limiter_id,json=rateLimiterId,proto3" json:"rate_limiter_id,omitempty"`
	Port          int64 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *RegisterNodeReq) Reset() {
	*x = RegisterNodeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_proxy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterNodeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterNodeReq) ProtoMessage() {}

func (x *RegisterNodeReq) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_proxy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterNodeReq.ProtoReflect.Descriptor instead.
func (*RegisterNodeReq) Descriptor() ([]byte, []int) {
	return file_proxy_proxy_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterNodeReq) GetRateLimiterId() int64 {
	if x != nil {
		return x.RateLimiterId
	}
	return 0
}

func (x *RegisterNodeReq) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

type RegisterNodeResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Res bool `protobuf:"varint,1,opt,name=res,proto3" json:"res,omitempty"`
}

func (x *RegisterNodeResp) Reset() {
	*x = RegisterNodeResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_proxy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterNodeResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterNodeResp) ProtoMessage() {}

func (x *RegisterNodeResp) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_proxy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterNodeResp.ProtoReflect.Descriptor instead.
func (*RegisterNodeResp) Descriptor() ([]byte, []int) {
	return file_proxy_proxy_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterNodeResp) GetRes() bool {
	if x != nil {
		return x.Res
	}
	return false
}

type AllowRequestReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiKey int64 `protobuf:"varint,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
}

func (x *AllowRequestReq) Reset() {
	*x = AllowRequestReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_proxy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllowRequestReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllowRequestReq) ProtoMessage() {}

func (x *AllowRequestReq) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_proxy_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllowRequestReq.ProtoReflect.Descriptor instead.
func (*AllowRequestReq) Descriptor() ([]byte, []int) {
	return file_proxy_proxy_proto_rawDescGZIP(), []int{2}
}

func (x *AllowRequestReq) GetApiKey() int64 {
	if x != nil {
		return x.ApiKey
	}
	return 0
}

type AllowRequestResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Res bool `protobuf:"varint,1,opt,name=res,proto3" json:"res,omitempty"`
}

func (x *AllowRequestResp) Reset() {
	*x = AllowRequestResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proxy_proxy_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllowRequestResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllowRequestResp) ProtoMessage() {}

func (x *AllowRequestResp) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_proxy_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllowRequestResp.ProtoReflect.Descriptor instead.
func (*AllowRequestResp) Descriptor() ([]byte, []int) {
	return file_proxy_proxy_proto_rawDescGZIP(), []int{3}
}

func (x *AllowRequestResp) GetRes() bool {
	if x != nil {
		return x.Res
	}
	return false
}

var File_proxy_proxy_proto protoreflect.FileDescriptor

var file_proxy_proxy_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x22, 0x4d, 0x0a, 0x0f, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x12, 0x26, 0x0a,
	0x0f, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x24, 0x0a, 0x10, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a,
	0x03, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x72, 0x65, 0x73, 0x22,
	0x2a, 0x0a, 0x0f, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x22, 0x24, 0x0a, 0x10, 0x41,
	0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12,
	0x10, 0x0a, 0x03, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x72, 0x65,
	0x73, 0x32, 0x8d, 0x01, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x12, 0x41, 0x0a, 0x0c, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x2e, 0x70, 0x72,
	0x6f, 0x78, 0x79, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x41,
	0x0a, 0x0c, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x41,
	0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22,
	0x00, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x68, 0x65, 0x2d, 0x6d, 0x61, 0x6e, 0x2d, 0x64, 0x61, 0x76, 0x69, 0x64, 0x2f, 0x64, 0x69, 0x73,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x2d, 0x72, 0x61, 0x74, 0x65, 0x2d, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x65, 0x72, 0x2d, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2f, 0x73, 0x72, 0x63, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proxy_proxy_proto_rawDescOnce sync.Once
	file_proxy_proxy_proto_rawDescData = file_proxy_proxy_proto_rawDesc
)

func file_proxy_proxy_proto_rawDescGZIP() []byte {
	file_proxy_proxy_proto_rawDescOnce.Do(func() {
		file_proxy_proxy_proto_rawDescData = protoimpl.X.CompressGZIP(file_proxy_proxy_proto_rawDescData)
	})
	return file_proxy_proxy_proto_rawDescData
}

var file_proxy_proxy_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proxy_proxy_proto_goTypes = []interface{}{
	(*RegisterNodeReq)(nil),  // 0: proxy.RegisterNodeReq
	(*RegisterNodeResp)(nil), // 1: proxy.RegisterNodeResp
	(*AllowRequestReq)(nil),  // 2: proxy.AllowRequestReq
	(*AllowRequestResp)(nil), // 3: proxy.AllowRequestResp
}
var file_proxy_proxy_proto_depIdxs = []int32{
	0, // 0: proxy.Proxy.RegisterNode:input_type -> proxy.RegisterNodeReq
	2, // 1: proxy.Proxy.AllowRequest:input_type -> proxy.AllowRequestReq
	1, // 2: proxy.Proxy.RegisterNode:output_type -> proxy.RegisterNodeResp
	3, // 3: proxy.Proxy.AllowRequest:output_type -> proxy.AllowRequestResp
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proxy_proxy_proto_init() }
func file_proxy_proxy_proto_init() {
	if File_proxy_proxy_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proxy_proxy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterNodeReq); i {
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
		file_proxy_proxy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterNodeResp); i {
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
		file_proxy_proxy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllowRequestReq); i {
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
		file_proxy_proxy_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllowRequestResp); i {
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
			RawDescriptor: file_proxy_proxy_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proxy_proxy_proto_goTypes,
		DependencyIndexes: file_proxy_proxy_proto_depIdxs,
		MessageInfos:      file_proxy_proxy_proto_msgTypes,
	}.Build()
	File_proxy_proxy_proto = out.File
	file_proxy_proxy_proto_rawDesc = nil
	file_proxy_proxy_proto_goTypes = nil
	file_proxy_proxy_proto_depIdxs = nil
}