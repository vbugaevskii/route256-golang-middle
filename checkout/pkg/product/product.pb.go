// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: product.proto

package product

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

type RequestGetProduct struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Sku   uint32 `protobuf:"varint,2,opt,name=sku,proto3" json:"sku,omitempty"`
}

func (x *RequestGetProduct) Reset() {
	*x = RequestGetProduct{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestGetProduct) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestGetProduct) ProtoMessage() {}

func (x *RequestGetProduct) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestGetProduct.ProtoReflect.Descriptor instead.
func (*RequestGetProduct) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{0}
}

func (x *RequestGetProduct) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *RequestGetProduct) GetSku() uint32 {
	if x != nil {
		return x.Sku
	}
	return 0
}

type ResponseGetProduct struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Price uint32 `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
}

func (x *ResponseGetProduct) Reset() {
	*x = ResponseGetProduct{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseGetProduct) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseGetProduct) ProtoMessage() {}

func (x *ResponseGetProduct) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseGetProduct.ProtoReflect.Descriptor instead.
func (*ResponseGetProduct) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseGetProduct) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ResponseGetProduct) GetPrice() uint32 {
	if x != nil {
		return x.Price
	}
	return 0
}

type RequestListSkus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token         string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	StartAfterSku uint32 `protobuf:"varint,2,opt,name=startAfterSku,proto3" json:"startAfterSku,omitempty"`
	Count         uint32 `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *RequestListSkus) Reset() {
	*x = RequestListSkus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestListSkus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestListSkus) ProtoMessage() {}

func (x *RequestListSkus) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestListSkus.ProtoReflect.Descriptor instead.
func (*RequestListSkus) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{2}
}

func (x *RequestListSkus) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *RequestListSkus) GetStartAfterSku() uint32 {
	if x != nil {
		return x.StartAfterSku
	}
	return 0
}

func (x *RequestListSkus) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ResponseListSkus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Skus []uint32 `protobuf:"varint,1,rep,packed,name=skus,proto3" json:"skus,omitempty"`
}

func (x *ResponseListSkus) Reset() {
	*x = ResponseListSkus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_product_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseListSkus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseListSkus) ProtoMessage() {}

func (x *ResponseListSkus) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseListSkus.ProtoReflect.Descriptor instead.
func (*ResponseListSkus) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{3}
}

func (x *ResponseListSkus) GetSkus() []uint32 {
	if x != nil {
		return x.Skus
	}
	return nil
}

var File_product_proto protoreflect.FileDescriptor

var file_product_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0x3b, 0x0a, 0x11, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x6b, 0x75, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x03, 0x73, 0x6b, 0x75, 0x22, 0x3e, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0x63, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x53, 0x6b, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x24,
	0x0a, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x66, 0x74, 0x65, 0x72, 0x53, 0x6b, 0x75, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x41, 0x66, 0x74, 0x65,
	0x72, 0x53, 0x6b, 0x75, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x26, 0x0a, 0x10, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6b, 0x75, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x6b, 0x75, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x04, 0x73, 0x6b,
	0x75, 0x73, 0x32, 0x8e, 0x01, 0x0a, 0x04, 0x4c, 0x6f, 0x6d, 0x73, 0x12, 0x45, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x47, 0x65, 0x74, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x12, 0x3f, 0x0a, 0x08, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6b, 0x75, 0x73, 0x12, 0x18,
	0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x53, 0x6b, 0x75, 0x73, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x53,
	0x6b, 0x75, 0x73, 0x42, 0x1e, 0x5a, 0x1c, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x32, 0x35, 0x36, 0x2f,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_product_proto_rawDescOnce sync.Once
	file_product_proto_rawDescData = file_product_proto_rawDesc
)

func file_product_proto_rawDescGZIP() []byte {
	file_product_proto_rawDescOnce.Do(func() {
		file_product_proto_rawDescData = protoimpl.X.CompressGZIP(file_product_proto_rawDescData)
	})
	return file_product_proto_rawDescData
}

var file_product_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_product_proto_goTypes = []interface{}{
	(*RequestGetProduct)(nil),  // 0: product.RequestGetProduct
	(*ResponseGetProduct)(nil), // 1: product.ResponseGetProduct
	(*RequestListSkus)(nil),    // 2: product.RequestListSkus
	(*ResponseListSkus)(nil),   // 3: product.ResponseListSkus
}
var file_product_proto_depIdxs = []int32{
	0, // 0: product.Loms.GetProduct:input_type -> product.RequestGetProduct
	2, // 1: product.Loms.ListSkus:input_type -> product.RequestListSkus
	1, // 2: product.Loms.GetProduct:output_type -> product.ResponseGetProduct
	3, // 3: product.Loms.ListSkus:output_type -> product.ResponseListSkus
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_product_proto_init() }
func file_product_proto_init() {
	if File_product_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_product_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestGetProduct); i {
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
		file_product_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseGetProduct); i {
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
		file_product_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestListSkus); i {
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
		file_product_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseListSkus); i {
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
			RawDescriptor: file_product_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_product_proto_goTypes,
		DependencyIndexes: file_product_proto_depIdxs,
		MessageInfos:      file_product_proto_msgTypes,
	}.Build()
	File_product_proto = out.File
	file_product_proto_rawDesc = nil
	file_product_proto_goTypes = nil
	file_product_proto_depIdxs = nil
}
