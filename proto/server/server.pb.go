// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

package server

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type PaymentRequestMessage struct {
	From                 string   `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To                   string   `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Amount               int64    `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PaymentRequestMessage) Reset()         { *m = PaymentRequestMessage{} }
func (m *PaymentRequestMessage) String() string { return proto.CompactTextString(m) }
func (*PaymentRequestMessage) ProtoMessage()    {}
func (*PaymentRequestMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{0}
}

func (m *PaymentRequestMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PaymentRequestMessage.Unmarshal(m, b)
}
func (m *PaymentRequestMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PaymentRequestMessage.Marshal(b, m, deterministic)
}
func (m *PaymentRequestMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PaymentRequestMessage.Merge(m, src)
}
func (m *PaymentRequestMessage) XXX_Size() int {
	return xxx_messageInfo_PaymentRequestMessage.Size(m)
}
func (m *PaymentRequestMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_PaymentRequestMessage.DiscardUnknown(m)
}

var xxx_messageInfo_PaymentRequestMessage proto.InternalMessageInfo

func (m *PaymentRequestMessage) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *PaymentRequestMessage) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *PaymentRequestMessage) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type Address struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{1}
}

func (m *Address) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Address.Unmarshal(m, b)
}
func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Address.Marshal(b, m, deterministic)
}
func (m *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(m, src)
}
func (m *Address) XXX_Size() int {
	return xxx_messageInfo_Address.Size(m)
}
func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type Result struct {
	Result               bool     `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{2}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

type CommunicationInfo struct {
	IPAddress            string   `protobuf:"bytes,1,opt,name=IPAddress,proto3" json:"IPAddress,omitempty"`
	Port                 int64    `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommunicationInfo) Reset()         { *m = CommunicationInfo{} }
func (m *CommunicationInfo) String() string { return proto.CompactTextString(m) }
func (*CommunicationInfo) ProtoMessage()    {}
func (*CommunicationInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{3}
}

func (m *CommunicationInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommunicationInfo.Unmarshal(m, b)
}
func (m *CommunicationInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommunicationInfo.Marshal(b, m, deterministic)
}
func (m *CommunicationInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommunicationInfo.Merge(m, src)
}
func (m *CommunicationInfo) XXX_Size() int {
	return xxx_messageInfo_CommunicationInfo.Size(m)
}
func (m *CommunicationInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_CommunicationInfo.DiscardUnknown(m)
}

var xxx_messageInfo_CommunicationInfo proto.InternalMessageInfo

func (m *CommunicationInfo) GetIPAddress() string {
	if m != nil {
		return m.IPAddress
	}
	return ""
}

func (m *CommunicationInfo) GetPort() int64 {
	if m != nil {
		return m.Port
	}
	return 0
}

func init() {
	proto.RegisterType((*PaymentRequestMessage)(nil), "paymentRequestMessage")
	proto.RegisterType((*Address)(nil), "address")
	proto.RegisterType((*Result)(nil), "Result")
	proto.RegisterType((*CommunicationInfo)(nil), "CommunicationInfo")
}

func init() { proto.RegisterFile("server.proto", fileDescriptor_ad098daeda4239f7) }

var fileDescriptor_ad098daeda4239f7 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x4d, 0x4b, 0x03, 0x31,
	0x10, 0x86, 0xbb, 0xbb, 0xb2, 0x6d, 0x07, 0x29, 0x38, 0x60, 0x59, 0x8a, 0xc2, 0x92, 0x53, 0x4f,
	0x39, 0xd8, 0x9b, 0x37, 0x11, 0x0f, 0x3d, 0x08, 0x92, 0xfe, 0x82, 0xd8, 0x4e, 0xa5, 0x60, 0x32,
	0x6b, 0x32, 0x2b, 0xf4, 0xdf, 0x4b, 0xd2, 0x2d, 0xe2, 0xc7, 0xed, 0x99, 0x09, 0x79, 0xf2, 0xbe,
	0x81, 0xcb, 0x48, 0xe1, 0x93, 0x82, 0xee, 0x02, 0x0b, 0xab, 0x0d, 0x5c, 0x77, 0xf6, 0xe8, 0xc8,
	0x8b, 0xa1, 0x8f, 0x9e, 0xa2, 0x3c, 0x53, 0x8c, 0xf6, 0x8d, 0x10, 0xe1, 0x62, 0x1f, 0xd8, 0x35,
	0x45, 0x5b, 0x2c, 0xa7, 0x26, 0x33, 0xce, 0xa0, 0x14, 0x6e, 0xca, 0xbc, 0x29, 0x85, 0x71, 0x0e,
	0xb5, 0x75, 0xdc, 0x7b, 0x69, 0xaa, 0xb6, 0x58, 0x56, 0x66, 0x98, 0xd4, 0x2d, 0x8c, 0xed, 0x6e,
	0x17, 0x28, 0xc6, 0xa4, 0x49, 0x78, 0xd6, 0x24, 0x56, 0x2d, 0xd4, 0x86, 0x62, 0xff, 0x2e, 0x49,
	0x10, 0x32, 0xe5, 0xf3, 0x89, 0x19, 0x26, 0xf5, 0x04, 0x57, 0x8f, 0xec, 0x5c, 0xef, 0x0f, 0x5b,
	0x2b, 0x07, 0xf6, 0x6b, 0xbf, 0x67, 0xbc, 0x81, 0xe9, 0xfa, 0xe5, 0xe1, 0xe4, 0x1d, 0x7c, 0xdf,
	0x8b, 0xf4, 0x50, 0xc7, 0x41, 0x72, 0xba, 0xca, 0x64, 0xbe, 0x3b, 0x42, 0xbd, 0xc9, 0x65, 0x71,
	0x05, 0xb3, 0x9f, 0x35, 0x71, 0xae, 0xff, 0xed, 0xbd, 0x18, 0xeb, 0x53, 0x36, 0x35, 0xc2, 0x7b,
	0x68, 0xb6, 0xbf, 0x53, 0x9c, 0xaf, 0x4f, 0xf4, 0xd0, 0x70, 0x81, 0xfa, 0x4f, 0x54, 0x35, 0x7a,
	0xad, 0xf3, 0xf7, 0xae, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x8d, 0x97, 0xd9, 0xe7, 0x6e, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerClient interface {
	PaymentRequest(ctx context.Context, in *PaymentRequestMessage, opts ...grpc.CallOption) (*Result, error)
	CommunicationInfoRequest(ctx context.Context, in *Address, opts ...grpc.CallOption) (*CommunicationInfo, error)
}

type serverClient struct {
	cc *grpc.ClientConn
}

func NewServerClient(cc *grpc.ClientConn) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) PaymentRequest(ctx context.Context, in *PaymentRequestMessage, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/Server/paymentRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) CommunicationInfoRequest(ctx context.Context, in *Address, opts ...grpc.CallOption) (*CommunicationInfo, error) {
	out := new(CommunicationInfo)
	err := c.cc.Invoke(ctx, "/Server/communicationInfoRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServer is the server API for Server service.
type ServerServer interface {
	PaymentRequest(context.Context, *PaymentRequestMessage) (*Result, error)
	CommunicationInfoRequest(context.Context, *Address) (*CommunicationInfo, error)
}

// UnimplementedServerServer can be embedded to have forward compatible implementations.
type UnimplementedServerServer struct {
}

func (*UnimplementedServerServer) PaymentRequest(ctx context.Context, req *PaymentRequestMessage) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PaymentRequest not implemented")
}
func (*UnimplementedServerServer) CommunicationInfoRequest(ctx context.Context, req *Address) (*CommunicationInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommunicationInfoRequest not implemented")
}

func RegisterServerServer(s *grpc.Server, srv ServerServer) {
	s.RegisterService(&_Server_serviceDesc, srv)
}

func _Server_PaymentRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaymentRequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).PaymentRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Server/PaymentRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).PaymentRequest(ctx, req.(*PaymentRequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_CommunicationInfoRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).CommunicationInfoRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Server/CommunicationInfoRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).CommunicationInfoRequest(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

var _Server_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "paymentRequest",
			Handler:    _Server_PaymentRequest_Handler,
		},
		{
			MethodName: "communicationInfoRequest",
			Handler:    _Server_CommunicationInfoRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
