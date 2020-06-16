// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Payload struct {
	TraceId              string            `protobuf:"bytes,1,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	Id                   string            `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	ServiceName          string            `protobuf:"bytes,3,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	ServiceVersion       string            `protobuf:"bytes,4,opt,name=service_version,json=serviceVersion,proto3" json:"service_version,omitempty"`
	Topic                string            `protobuf:"bytes,5,opt,name=topic,proto3" json:"topic,omitempty"`
	Header               map[string]string `protobuf:"bytes,6,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ContentType          string            `protobuf:"bytes,7,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	Body                 []byte            `protobuf:"bytes,8,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Payload) Reset()         { *m = Payload{} }
func (m *Payload) String() string { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()    {}
func (*Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *Payload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Payload.Unmarshal(m, b)
}
func (m *Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Payload.Marshal(b, m, deterministic)
}
func (m *Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload.Merge(m, src)
}
func (m *Payload) XXX_Size() int {
	return xxx_messageInfo_Payload.Size(m)
}
func (m *Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_Payload proto.InternalMessageInfo

func (m *Payload) GetTraceId() string {
	if m != nil {
		return m.TraceId
	}
	return ""
}

func (m *Payload) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Payload) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *Payload) GetServiceVersion() string {
	if m != nil {
		return m.ServiceVersion
	}
	return ""
}

func (m *Payload) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *Payload) GetHeader() map[string]string {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Payload) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *Payload) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterType((*Payload)(nil), "proto.Payload")
	proto.RegisterMapType((map[string]string)(nil), "proto.Payload.HeaderEntry")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0x69, 0xbb, 0x6d, 0xd7, 0xe9, 0xba, 0xca, 0xe0, 0x21, 0xee, 0xa9, 0x7a, 0xb1, 0xa7,
	0x1e, 0xd6, 0x8b, 0x7a, 0x17, 0xf4, 0x22, 0x52, 0xc4, 0x6b, 0xc9, 0x36, 0x83, 0x06, 0xb7, 0x4d,
	0x49, 0x63, 0x21, 0xef, 0xe6, 0xc3, 0x49, 0x93, 0x08, 0x9e, 0x32, 0xf3, 0xfd, 0x3f, 0xe1, 0x63,
	0xe0, 0xb4, 0xa7, 0x69, 0xe2, 0x1f, 0x54, 0x8f, 0x5a, 0x19, 0x85, 0xa9, 0x7b, 0xae, 0x7f, 0x62,
	0xc8, 0x5f, 0xb9, 0x3d, 0x2a, 0x2e, 0xf0, 0x12, 0xd6, 0x46, 0xf3, 0x8e, 0x5a, 0x29, 0x58, 0x54,
	0x46, 0xd5, 0x49, 0x93, 0xbb, 0xfd, 0x59, 0xe0, 0x16, 0x62, 0x29, 0x58, 0xec, 0x60, 0x2c, 0x05,
	0x5e, 0xc1, 0x66, 0x22, 0x3d, 0xcb, 0x8e, 0xda, 0x81, 0xf7, 0xc4, 0x12, 0x97, 0x14, 0x81, 0xbd,
	0xf0, 0x9e, 0xf0, 0x06, 0xce, 0xfe, 0x2a, 0x33, 0xe9, 0x49, 0xaa, 0x81, 0xad, 0x5c, 0x6b, 0x1b,
	0xf0, 0xbb, 0xa7, 0x78, 0x01, 0xa9, 0x51, 0xa3, 0xec, 0x58, 0xea, 0x62, 0xbf, 0xe0, 0x1e, 0xb2,
	0x4f, 0xe2, 0x82, 0x34, 0xcb, 0xca, 0xa4, 0x2a, 0xf6, 0x3b, 0xef, 0x5d, 0x07, 0xd9, 0xfa, 0xc9,
	0x85, 0x8f, 0x83, 0xd1, 0xb6, 0x09, 0xcd, 0xc5, 0xaa, 0x53, 0x83, 0xa1, 0xc1, 0xb4, 0xc6, 0x8e,
	0xc4, 0x72, 0x6f, 0x15, 0xd8, 0x9b, 0x1d, 0x09, 0x11, 0x56, 0x07, 0x25, 0x2c, 0x5b, 0x97, 0x51,
	0xb5, 0x69, 0xdc, 0xbc, 0xbb, 0x87, 0xe2, 0xdf, 0x6f, 0x78, 0x0e, 0xc9, 0x17, 0xd9, 0x70, 0x81,
	0x65, 0x5c, 0x0c, 0x67, 0x7e, 0xfc, 0xa6, 0x70, 0x00, 0xbf, 0x3c, 0xc4, 0x77, 0xd1, 0x21, 0x73,
	0x52, 0xb7, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x85, 0x90, 0xe0, 0xd6, 0x5d, 0x01, 0x00, 0x00,
}
