// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.4
// source: internal/proto/metric.proto

package internal

import (
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/any"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type DatasetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetricId   string               `protobuf:"bytes,1,opt,name=MetricId,proto3" json:"MetricId,omitempty"`
	UserId     string               `protobuf:"bytes,2,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Timestamp  *timestamp.Timestamp `protobuf:"bytes,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	So         string               `protobuf:"bytes,4,opt,name=So,proto3" json:"So,omitempty"`
	RitVersion string               `protobuf:"bytes,5,opt,name=RitVersion,proto3" json:"RitVersion,omitempty"`
	Data       []byte               `protobuf:"bytes,6,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *DatasetRequest) Reset() {
	*x = DatasetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatasetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatasetRequest) ProtoMessage() {}

func (x *DatasetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatasetRequest.ProtoReflect.Descriptor instead.
func (*DatasetRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_metric_proto_rawDescGZIP(), []int{0}
}

func (x *DatasetRequest) GetMetricId() string {
	if x != nil {
		return x.MetricId
	}
	return ""
}

func (x *DatasetRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DatasetRequest) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *DatasetRequest) GetSo() string {
	if x != nil {
		return x.So
	}
	return ""
}

func (x *DatasetRequest) GetRitVersion() string {
	if x != nil {
		return x.RitVersion
	}
	return ""
}

func (x *DatasetRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_internal_proto_metric_proto protoreflect.FileDescriptor

var file_internal_proto_metric_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xc2, 0x01, 0x0a, 0x0e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x0e, 0x0a, 0x02, 0x53, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x53,
	0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x69, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x52, 0x69, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x32, 0x4a, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x12, 0x3d, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x2e,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_metric_proto_rawDescOnce sync.Once
	file_internal_proto_metric_proto_rawDescData = file_internal_proto_metric_proto_rawDesc
)

func file_internal_proto_metric_proto_rawDescGZIP() []byte {
	file_internal_proto_metric_proto_rawDescOnce.Do(func() {
		file_internal_proto_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_metric_proto_rawDescData)
	})
	return file_internal_proto_metric_proto_rawDescData
}

var file_internal_proto_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_internal_proto_metric_proto_goTypes = []interface{}{
	(*DatasetRequest)(nil),      // 0: internal.DatasetRequest
	(*timestamp.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*empty.Empty)(nil),         // 2: google.protobuf.Empty
}
var file_internal_proto_metric_proto_depIdxs = []int32{
	1, // 0: internal.DatasetRequest.Timestamp:type_name -> google.protobuf.Timestamp
	0, // 1: internal.Processor.Process:input_type -> internal.DatasetRequest
	2, // 2: internal.Processor.Process:output_type -> google.protobuf.Empty
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_proto_metric_proto_init() }
func file_internal_proto_metric_proto_init() {
	if File_internal_proto_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatasetRequest); i {
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
			RawDescriptor: file_internal_proto_metric_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_metric_proto_goTypes,
		DependencyIndexes: file_internal_proto_metric_proto_depIdxs,
		MessageInfos:      file_internal_proto_metric_proto_msgTypes,
	}.Build()
	File_internal_proto_metric_proto = out.File
	file_internal_proto_metric_proto_rawDesc = nil
	file_internal_proto_metric_proto_goTypes = nil
	file_internal_proto_metric_proto_depIdxs = nil
}
