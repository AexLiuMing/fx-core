// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/staking/v1beta1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgGrantPrivilege defines the GrantPrivilege message.
type MsgGrantPrivilege struct {
	ValidatorAddress string     `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	FromAddress      string     `protobuf:"bytes,2,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	ToPubkey         *types.Any `protobuf:"bytes,3,opt,name=to_pubkey,json=toPubkey,proto3" json:"to_pubkey,omitempty"`
	Signature        string     `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *MsgGrantPrivilege) Reset()         { *m = MsgGrantPrivilege{} }
func (m *MsgGrantPrivilege) String() string { return proto.CompactTextString(m) }
func (*MsgGrantPrivilege) ProtoMessage()    {}
func (*MsgGrantPrivilege) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c6db236fc9c571c, []int{0}
}
func (m *MsgGrantPrivilege) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgGrantPrivilege) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgGrantPrivilege.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgGrantPrivilege) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgGrantPrivilege.Merge(m, src)
}
func (m *MsgGrantPrivilege) XXX_Size() int {
	return m.Size()
}
func (m *MsgGrantPrivilege) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgGrantPrivilege.DiscardUnknown(m)
}

var xxx_messageInfo_MsgGrantPrivilege proto.InternalMessageInfo

func (m *MsgGrantPrivilege) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *MsgGrantPrivilege) GetFromAddress() string {
	if m != nil {
		return m.FromAddress
	}
	return ""
}

func (m *MsgGrantPrivilege) GetToPubkey() *types.Any {
	if m != nil {
		return m.ToPubkey
	}
	return nil
}

func (m *MsgGrantPrivilege) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

type MsgGrantPrivilegeResponse struct {
}

func (m *MsgGrantPrivilegeResponse) Reset()         { *m = MsgGrantPrivilegeResponse{} }
func (m *MsgGrantPrivilegeResponse) String() string { return proto.CompactTextString(m) }
func (*MsgGrantPrivilegeResponse) ProtoMessage()    {}
func (*MsgGrantPrivilegeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c6db236fc9c571c, []int{1}
}
func (m *MsgGrantPrivilegeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgGrantPrivilegeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgGrantPrivilegeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgGrantPrivilegeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgGrantPrivilegeResponse.Merge(m, src)
}
func (m *MsgGrantPrivilegeResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgGrantPrivilegeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgGrantPrivilegeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgGrantPrivilegeResponse proto.InternalMessageInfo

type MsgEditConsensusPubKey struct {
	ValidatorAddress string     `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	From             string     `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Pubkey           *types.Any `protobuf:"bytes,3,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
}

func (m *MsgEditConsensusPubKey) Reset()         { *m = MsgEditConsensusPubKey{} }
func (m *MsgEditConsensusPubKey) String() string { return proto.CompactTextString(m) }
func (*MsgEditConsensusPubKey) ProtoMessage()    {}
func (*MsgEditConsensusPubKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c6db236fc9c571c, []int{2}
}
func (m *MsgEditConsensusPubKey) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgEditConsensusPubKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgEditConsensusPubKey.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgEditConsensusPubKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEditConsensusPubKey.Merge(m, src)
}
func (m *MsgEditConsensusPubKey) XXX_Size() int {
	return m.Size()
}
func (m *MsgEditConsensusPubKey) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEditConsensusPubKey.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEditConsensusPubKey proto.InternalMessageInfo

func (m *MsgEditConsensusPubKey) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *MsgEditConsensusPubKey) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *MsgEditConsensusPubKey) GetPubkey() *types.Any {
	if m != nil {
		return m.Pubkey
	}
	return nil
}

type MsgEditConsensusPubKeyResponse struct {
}

func (m *MsgEditConsensusPubKeyResponse) Reset()         { *m = MsgEditConsensusPubKeyResponse{} }
func (m *MsgEditConsensusPubKeyResponse) String() string { return proto.CompactTextString(m) }
func (*MsgEditConsensusPubKeyResponse) ProtoMessage()    {}
func (*MsgEditConsensusPubKeyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c6db236fc9c571c, []int{3}
}
func (m *MsgEditConsensusPubKeyResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgEditConsensusPubKeyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgEditConsensusPubKeyResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgEditConsensusPubKeyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEditConsensusPubKeyResponse.Merge(m, src)
}
func (m *MsgEditConsensusPubKeyResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgEditConsensusPubKeyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEditConsensusPubKeyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEditConsensusPubKeyResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgGrantPrivilege)(nil), "fx.staking.v1.MsgGrantPrivilege")
	proto.RegisterType((*MsgGrantPrivilegeResponse)(nil), "fx.staking.v1.MsgGrantPrivilegeResponse")
	proto.RegisterType((*MsgEditConsensusPubKey)(nil), "fx.staking.v1.MsgEditConsensusPubKey")
	proto.RegisterType((*MsgEditConsensusPubKeyResponse)(nil), "fx.staking.v1.MsgEditConsensusPubKeyResponse")
}

func init() { proto.RegisterFile("fx/staking/v1beta1/tx.proto", fileDescriptor_1c6db236fc9c571c) }

var fileDescriptor_1c6db236fc9c571c = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xdd, 0x8a, 0xd3, 0x40,
	0x14, 0xee, 0xb8, 0xcb, 0x62, 0x67, 0x55, 0xdc, 0x71, 0x91, 0x34, 0x2b, 0x21, 0x06, 0x84, 0x82,
	0x76, 0x86, 0xae, 0x4f, 0xb0, 0x2b, 0xae, 0x17, 0x4b, 0xa1, 0xf4, 0x52, 0x84, 0x32, 0x49, 0x26,
	0xe3, 0xd0, 0x76, 0x26, 0xcc, 0x4c, 0x42, 0xf2, 0x16, 0xbe, 0x82, 0xef, 0xe0, 0x43, 0x88, 0x17,
	0xd2, 0x4b, 0x2f, 0xa5, 0x7d, 0x11, 0x69, 0x7e, 0x5a, 0xb4, 0x05, 0x95, 0xbd, 0xcb, 0xf9, 0xbe,
	0xef, 0x9c, 0xf3, 0x7d, 0x27, 0x09, 0xbc, 0x48, 0x0a, 0x62, 0x2c, 0x9d, 0x09, 0xc9, 0x49, 0x3e,
	0x0c, 0x99, 0xa5, 0x43, 0x62, 0x0b, 0x9c, 0x6a, 0x65, 0x15, 0x7a, 0x98, 0x14, 0xb8, 0x21, 0x71,
	0x3e, 0x74, 0x7b, 0x5c, 0x29, 0x3e, 0x67, 0xa4, 0x22, 0xc3, 0x2c, 0x21, 0x54, 0x96, 0xb5, 0xd2,
	0xed, 0x45, 0xca, 0x2c, 0x94, 0x99, 0x56, 0x15, 0xa9, 0x8b, 0x9a, 0x0a, 0xbe, 0x03, 0x78, 0x36,
	0x32, 0xfc, 0x9d, 0xa6, 0xd2, 0x8e, 0xb5, 0xc8, 0xc5, 0x9c, 0x71, 0x86, 0x5e, 0xc2, 0xb3, 0x9c,
	0xce, 0x45, 0x4c, 0xad, 0xd2, 0x53, 0x1a, 0xc7, 0x9a, 0x19, 0xe3, 0x00, 0x1f, 0xf4, 0xbb, 0x93,
	0xc7, 0x5b, 0xe2, 0xaa, 0xc6, 0xd1, 0x73, 0xf8, 0x20, 0xd1, 0x6a, 0xb1, 0xd5, 0xdd, 0xab, 0x74,
	0xa7, 0x1b, 0xac, 0x95, 0xdc, 0xc2, 0xae, 0x55, 0xd3, 0x34, 0x0b, 0x67, 0xac, 0x74, 0x8e, 0x7c,
	0xd0, 0x3f, 0xbd, 0x3c, 0xc7, 0xb5, 0x5f, 0xdc, 0xfa, 0xc5, 0x57, 0xb2, 0xbc, 0x76, 0xbe, 0x7d,
	0x19, 0x9c, 0x37, 0x06, 0x23, 0x5d, 0xa6, 0x56, 0xe1, 0x71, 0x16, 0xde, 0xb2, 0x72, 0x72, 0xdf,
	0xaa, 0x71, 0xd5, 0x8f, 0x9e, 0xc1, 0xae, 0x11, 0x5c, 0x52, 0x9b, 0x69, 0xe6, 0x1c, 0x57, 0xcb,
	0x76, 0x40, 0x70, 0x01, 0x7b, 0x7b, 0x79, 0x26, 0xcc, 0xa4, 0x4a, 0x1a, 0x16, 0x7c, 0x06, 0xf0,
	0xe9, 0xc8, 0xf0, 0xb7, 0xb1, 0xb0, 0x6f, 0x36, 0x80, 0x34, 0x99, 0xa9, 0xe7, 0xff, 0x5f, 0x64,
	0x04, 0x8f, 0x37, 0xf1, 0x9a, 0xa8, 0xd5, 0x33, 0xba, 0x81, 0x27, 0x77, 0x0a, 0xd8, 0x74, 0x07,
	0x3e, 0xf4, 0x0e, 0x5b, 0x6c, 0x53, 0x5c, 0x2e, 0x01, 0x3c, 0x1a, 0x19, 0x8e, 0x3e, 0xc0, 0x47,
	0x7f, 0xbc, 0x37, 0x1f, 0xff, 0xf6, 0x4d, 0xe0, 0xbd, 0x4b, 0xb8, 0xfd, 0xbf, 0x29, 0xda, 0x2d,
	0x68, 0x06, 0x9f, 0x1c, 0xba, 0xd3, 0x8b, 0xfd, 0x01, 0x07, 0x64, 0xee, 0xe0, 0x9f, 0x64, 0xed,
	0xb2, 0xeb, 0x9b, 0xaf, 0x2b, 0x0f, 0x2c, 0x57, 0x1e, 0xf8, 0xb9, 0xf2, 0xc0, 0xa7, 0xb5, 0xd7,
	0x59, 0xae, 0xbd, 0xce, 0x8f, 0xb5, 0xd7, 0x79, 0xff, 0x8a, 0x0b, 0xfb, 0x31, 0x0b, 0x71, 0xa4,
	0x16, 0x24, 0xc9, 0x64, 0x64, 0x85, 0x92, 0x05, 0x49, 0x8a, 0x41, 0xa4, 0x34, 0x23, 0xbb, 0xdf,
	0xc3, 0x96, 0x29, 0x33, 0xe1, 0x49, 0x75, 0xec, 0xd7, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x48,
	0x74, 0x03, 0x05, 0x39, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// GrantPrivilege defines a method for granting privilege to a validator.
	GrantPrivilege(ctx context.Context, in *MsgGrantPrivilege, opts ...grpc.CallOption) (*MsgGrantPrivilegeResponse, error)
	// EditConsensusKey defines a method for editing consensus pubkey of a validator.
	EditConsensusPubKey(ctx context.Context, in *MsgEditConsensusPubKey, opts ...grpc.CallOption) (*MsgEditConsensusPubKeyResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) GrantPrivilege(ctx context.Context, in *MsgGrantPrivilege, opts ...grpc.CallOption) (*MsgGrantPrivilegeResponse, error) {
	out := new(MsgGrantPrivilegeResponse)
	err := c.cc.Invoke(ctx, "/fx.staking.v1.Msg/GrantPrivilege", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) EditConsensusPubKey(ctx context.Context, in *MsgEditConsensusPubKey, opts ...grpc.CallOption) (*MsgEditConsensusPubKeyResponse, error) {
	out := new(MsgEditConsensusPubKeyResponse)
	err := c.cc.Invoke(ctx, "/fx.staking.v1.Msg/EditConsensusPubKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// GrantPrivilege defines a method for granting privilege to a validator.
	GrantPrivilege(context.Context, *MsgGrantPrivilege) (*MsgGrantPrivilegeResponse, error)
	// EditConsensusKey defines a method for editing consensus pubkey of a validator.
	EditConsensusPubKey(context.Context, *MsgEditConsensusPubKey) (*MsgEditConsensusPubKeyResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) GrantPrivilege(ctx context.Context, req *MsgGrantPrivilege) (*MsgGrantPrivilegeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrantPrivilege not implemented")
}
func (*UnimplementedMsgServer) EditConsensusPubKey(ctx context.Context, req *MsgEditConsensusPubKey) (*MsgEditConsensusPubKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditConsensusPubKey not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_GrantPrivilege_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgGrantPrivilege)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).GrantPrivilege(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.staking.v1.Msg/GrantPrivilege",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).GrantPrivilege(ctx, req.(*MsgGrantPrivilege))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_EditConsensusPubKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgEditConsensusPubKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).EditConsensusPubKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.staking.v1.Msg/EditConsensusPubKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).EditConsensusPubKey(ctx, req.(*MsgEditConsensusPubKey))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fx.staking.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GrantPrivilege",
			Handler:    _Msg_GrantPrivilege_Handler,
		},
		{
			MethodName: "EditConsensusPubKey",
			Handler:    _Msg_EditConsensusPubKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fx/staking/v1beta1/tx.proto",
}

func (m *MsgGrantPrivilege) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgGrantPrivilege) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgGrantPrivilege) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x22
	}
	if m.ToPubkey != nil {
		{
			size, err := m.ToPubkey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.FromAddress) > 0 {
		i -= len(m.FromAddress)
		copy(dAtA[i:], m.FromAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.FromAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgGrantPrivilegeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgGrantPrivilegeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgGrantPrivilegeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgEditConsensusPubKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgEditConsensusPubKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgEditConsensusPubKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pubkey != nil {
		{
			size, err := m.Pubkey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.From) > 0 {
		i -= len(m.From)
		copy(dAtA[i:], m.From)
		i = encodeVarintTx(dAtA, i, uint64(len(m.From)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgEditConsensusPubKeyResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgEditConsensusPubKeyResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgEditConsensusPubKeyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgGrantPrivilege) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.FromAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.ToPubkey != nil {
		l = m.ToPubkey.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgGrantPrivilegeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgEditConsensusPubKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Pubkey != nil {
		l = m.Pubkey.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgEditConsensusPubKeyResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgGrantPrivilege) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgGrantPrivilege: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgGrantPrivilege: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FromAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToPubkey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ToPubkey == nil {
				m.ToPubkey = &types.Any{}
			}
			if err := m.ToPubkey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgGrantPrivilegeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgGrantPrivilegeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgGrantPrivilegeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgEditConsensusPubKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgEditConsensusPubKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgEditConsensusPubKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pubkey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pubkey == nil {
				m.Pubkey = &types.Any{}
			}
			if err := m.Pubkey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgEditConsensusPubKeyResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgEditConsensusPubKeyResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgEditConsensusPubKeyResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
