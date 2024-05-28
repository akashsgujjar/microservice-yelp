// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: proto/detail/detail.proto

package detail

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

type PostDetailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RestaurantName string `protobuf:"bytes,1,opt,name=restaurant_name,json=restaurantName,proto3" json:"restaurant_name,omitempty"`
	Location       string `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Style          string `protobuf:"bytes,3,opt,name=style,proto3" json:"style,omitempty"`
	Capacity       int32  `protobuf:"varint,4,opt,name=capacity,proto3" json:"capacity,omitempty"`
}

func (x *PostDetailRequest) Reset() {
	*x = PostDetailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_detail_detail_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostDetailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostDetailRequest) ProtoMessage() {}

func (x *PostDetailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_detail_detail_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostDetailRequest.ProtoReflect.Descriptor instead.
func (*PostDetailRequest) Descriptor() ([]byte, []int) {
	return file_proto_detail_detail_proto_rawDescGZIP(), []int{0}
}

func (x *PostDetailRequest) GetRestaurantName() string {
	if x != nil {
		return x.RestaurantName
	}
	return ""
}

func (x *PostDetailRequest) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *PostDetailRequest) GetStyle() string {
	if x != nil {
		return x.Style
	}
	return ""
}

func (x *PostDetailRequest) GetCapacity() int32 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

type PostDetailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *PostDetailResponse) Reset() {
	*x = PostDetailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_detail_detail_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostDetailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostDetailResponse) ProtoMessage() {}

func (x *PostDetailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_detail_detail_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostDetailResponse.ProtoReflect.Descriptor instead.
func (*PostDetailResponse) Descriptor() ([]byte, []int) {
	return file_proto_detail_detail_proto_rawDescGZIP(), []int{1}
}

func (x *PostDetailResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

// The request message containing the restaurant name
type GetDetailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RestaurantName string `protobuf:"bytes,1,opt,name=restaurant_name,json=restaurantName,proto3" json:"restaurant_name,omitempty"`
}

func (x *GetDetailRequest) Reset() {
	*x = GetDetailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_detail_detail_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDetailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDetailRequest) ProtoMessage() {}

func (x *GetDetailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_detail_detail_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDetailRequest.ProtoReflect.Descriptor instead.
func (*GetDetailRequest) Descriptor() ([]byte, []int) {
	return file_proto_detail_detail_proto_rawDescGZIP(), []int{2}
}

func (x *GetDetailRequest) GetRestaurantName() string {
	if x != nil {
		return x.RestaurantName
	}
	return ""
}

// The response message containing the restaurant details
type GetDetailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RestaurantName string `protobuf:"bytes,1,opt,name=restaurant_name,json=restaurantName,proto3" json:"restaurant_name,omitempty"`
	Location       string `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Style          string `protobuf:"bytes,3,opt,name=style,proto3" json:"style,omitempty"`
	Capacity       int32  `protobuf:"varint,4,opt,name=capacity,proto3" json:"capacity,omitempty"`
}

func (x *GetDetailResponse) Reset() {
	*x = GetDetailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_detail_detail_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDetailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDetailResponse) ProtoMessage() {}

func (x *GetDetailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_detail_detail_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDetailResponse.ProtoReflect.Descriptor instead.
func (*GetDetailResponse) Descriptor() ([]byte, []int) {
	return file_proto_detail_detail_proto_rawDescGZIP(), []int{3}
}

func (x *GetDetailResponse) GetRestaurantName() string {
	if x != nil {
		return x.RestaurantName
	}
	return ""
}

func (x *GetDetailResponse) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *GetDetailResponse) GetStyle() string {
	if x != nil {
		return x.Style
	}
	return ""
}

func (x *GetDetailResponse) GetCapacity() int32 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

var File_proto_detail_detail_proto protoreflect.FileDescriptor

var file_proto_detail_detail_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x2f, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x22, 0x8a, 0x01, 0x0a, 0x11, 0x50, 0x6f, 0x73, 0x74, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x74, 0x79, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79,
	0x22, 0x2c, 0x0a, 0x12, 0x50, 0x6f, 0x73, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x3b,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x8a, 0x01, 0x0a, 0x11,
	0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x73, 0x74,
	0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x32, 0x96, 0x01, 0x0a, 0x0d, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0a, 0x50, 0x6f,
	0x73, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x19, 0x2e, 0x64, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x2e, 0x50, 0x6f, 0x73,
	0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x40, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x2e, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x2e,
	0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_detail_detail_proto_rawDescOnce sync.Once
	file_proto_detail_detail_proto_rawDescData = file_proto_detail_detail_proto_rawDesc
)

func file_proto_detail_detail_proto_rawDescGZIP() []byte {
	file_proto_detail_detail_proto_rawDescOnce.Do(func() {
		file_proto_detail_detail_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_detail_detail_proto_rawDescData)
	})
	return file_proto_detail_detail_proto_rawDescData
}

var file_proto_detail_detail_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_detail_detail_proto_goTypes = []interface{}{
	(*PostDetailRequest)(nil),  // 0: detail.PostDetailRequest
	(*PostDetailResponse)(nil), // 1: detail.PostDetailResponse
	(*GetDetailRequest)(nil),   // 2: detail.GetDetailRequest
	(*GetDetailResponse)(nil),  // 3: detail.GetDetailResponse
}
var file_proto_detail_detail_proto_depIdxs = []int32{
	0, // 0: detail.DetailService.PostDetail:input_type -> detail.PostDetailRequest
	2, // 1: detail.DetailService.GetDetail:input_type -> detail.GetDetailRequest
	1, // 2: detail.DetailService.PostDetail:output_type -> detail.PostDetailResponse
	3, // 3: detail.DetailService.GetDetail:output_type -> detail.GetDetailResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_detail_detail_proto_init() }
func file_proto_detail_detail_proto_init() {
	if File_proto_detail_detail_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_detail_detail_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostDetailRequest); i {
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
		file_proto_detail_detail_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostDetailResponse); i {
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
		file_proto_detail_detail_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDetailRequest); i {
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
		file_proto_detail_detail_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDetailResponse); i {
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
			RawDescriptor: file_proto_detail_detail_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_detail_detail_proto_goTypes,
		DependencyIndexes: file_proto_detail_detail_proto_depIdxs,
		MessageInfos:      file_proto_detail_detail_proto_msgTypes,
	}.Build()
	File_proto_detail_detail_proto = out.File
	file_proto_detail_detail_proto_rawDesc = nil
	file_proto_detail_detail_proto_goTypes = nil
	file_proto_detail_detail_proto_depIdxs = nil
}
