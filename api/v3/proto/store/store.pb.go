// Code generated by protoc-gen-go. DO NOT EDIT.
// source: store.proto

package store

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_store_2bf46e8be68952db, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "store.Message")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TemporalStoreClient is the client API for TemporalStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TemporalStoreClient interface {
	Status(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
}

type temporalStoreClient struct {
	cc *grpc.ClientConn
}

func NewTemporalStoreClient(cc *grpc.ClientConn) TemporalStoreClient {
	return &temporalStoreClient{cc}
}

func (c *temporalStoreClient) Status(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/store.TemporalStore/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TemporalStoreServer is the server API for TemporalStore service.
type TemporalStoreServer interface {
	Status(context.Context, *Message) (*Message, error)
}

func RegisterTemporalStoreServer(s *grpc.Server, srv TemporalStoreServer) {
	s.RegisterService(&_TemporalStore_serviceDesc, srv)
}

func _TemporalStore_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TemporalStoreServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/store.TemporalStore/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TemporalStoreServer).Status(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _TemporalStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "store.TemporalStore",
	HandlerType: (*TemporalStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _TemporalStore_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "store.proto",
}

func init() { proto.RegisterFile("store.proto", fileDescriptor_store_2bf46e8be68952db) }

var fileDescriptor_store_2bf46e8be68952db = []byte{
	// 169 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x2e, 0xc9, 0x2f,
	0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0xa4, 0x64, 0xd2, 0xf3, 0xf3,
	0xd3, 0x73, 0x52, 0xf5, 0x13, 0x0b, 0x32, 0xf5, 0x13, 0xf3, 0xf2, 0xf2, 0x4b, 0x12, 0x4b, 0x32,
	0xf3, 0xf3, 0x8a, 0x21, 0x8a, 0x94, 0x94, 0xb9, 0xd8, 0x7d, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53,
	0x85, 0x24, 0xb8, 0xd8, 0x73, 0x21, 0x4c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0xd7,
	0x28, 0x98, 0x8b, 0x37, 0x24, 0x35, 0xb7, 0x20, 0xbf, 0x28, 0x31, 0x27, 0x18, 0x64, 0xa6, 0x90,
	0x13, 0x17, 0x5b, 0x70, 0x49, 0x62, 0x49, 0x69, 0xb1, 0x10, 0x9f, 0x1e, 0xc4, 0x4a, 0xa8, 0x21,
	0x52, 0x68, 0x7c, 0x25, 0x89, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0x09, 0x09, 0x09, 0xe8, 0x97, 0x19,
	0xeb, 0x83, 0xa5, 0xf4, 0x8b, 0xc1, 0x3a, 0x9d, 0xc4, 0xb8, 0x44, 0x92, 0x73, 0xf2, 0x4b, 0x53,
	0xf4, 0x4a, 0xa0, 0x46, 0x43, 0x74, 0x26, 0xb1, 0x81, 0x1d, 0x66, 0x0c, 0x08, 0x00, 0x00, 0xff,
	0xff, 0xcd, 0xcd, 0x8b, 0x48, 0xcc, 0x00, 0x00, 0x00,
}
