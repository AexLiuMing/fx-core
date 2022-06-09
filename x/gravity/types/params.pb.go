// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gravity/v1/params.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type Params struct {
	GravityId                      string                                 `protobuf:"bytes,1,opt,name=gravity_id,json=gravityId,proto3" json:"gravity_id,omitempty"`
	ContractSourceHash             string                                 `protobuf:"bytes,2,opt,name=contract_source_hash,json=contractSourceHash,proto3" json:"contract_source_hash,omitempty"`
	BridgeEthAddress               string                                 `protobuf:"bytes,4,opt,name=bridge_eth_address,json=bridgeEthAddress,proto3" json:"bridge_eth_address,omitempty"`
	BridgeChainId                  uint64                                 `protobuf:"varint,5,opt,name=bridge_chain_id,json=bridgeChainId,proto3" json:"bridge_chain_id,omitempty"`
	SignedValsetsWindow            uint64                                 `protobuf:"varint,6,opt,name=signed_valsets_window,json=signedValsetsWindow,proto3" json:"signed_valsets_window,omitempty"`
	SignedBatchesWindow            uint64                                 `protobuf:"varint,7,opt,name=signed_batches_window,json=signedBatchesWindow,proto3" json:"signed_batches_window,omitempty"`
	SignedClaimsWindow             uint64                                 `protobuf:"varint,8,opt,name=signed_claims_window,json=signedClaimsWindow,proto3" json:"signed_claims_window,omitempty"`
	TargetBatchTimeout             uint64                                 `protobuf:"varint,10,opt,name=target_batch_timeout,json=targetBatchTimeout,proto3" json:"target_batch_timeout,omitempty"`
	AverageBlockTime               uint64                                 `protobuf:"varint,11,opt,name=average_block_time,json=averageBlockTime,proto3" json:"average_block_time,omitempty"`
	AverageEthBlockTime            uint64                                 `protobuf:"varint,12,opt,name=average_eth_block_time,json=averageEthBlockTime,proto3" json:"average_eth_block_time,omitempty"`
	SlashFractionValset            github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,13,opt,name=slash_fraction_valset,json=slashFractionValset,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_valset"`
	SlashFractionBatch             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,14,opt,name=slash_fraction_batch,json=slashFractionBatch,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_batch"`
	SlashFractionClaim             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,15,opt,name=slash_fraction_claim,json=slashFractionClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_claim"`
	SlashFractionConflictingClaim  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,16,opt,name=slash_fraction_conflicting_claim,json=slashFractionConflictingClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_conflicting_claim"`
	UnbondSlashingValsetsWindow    uint64                                 `protobuf:"varint,17,opt,name=unbond_slashing_valsets_window,json=unbondSlashingValsetsWindow,proto3" json:"unbond_slashing_valsets_window,omitempty"`
	IbcTransferTimeoutHeight       uint64                                 `protobuf:"varint,18,opt,name=ibc_transfer_timeout_height,json=ibcTransferTimeoutHeight,proto3" json:"ibc_transfer_timeout_height,omitempty"`
	ValsetUpdatePowerChangePercent github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,19,opt,name=valset_update_power_change_percent,json=valsetUpdatePowerChangePercent,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"valset_update_power_change_percent"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_8772bac9489530eb, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetGravityId() string {
	if m != nil {
		return m.GravityId
	}
	return ""
}

func (m *Params) GetContractSourceHash() string {
	if m != nil {
		return m.ContractSourceHash
	}
	return ""
}

func (m *Params) GetBridgeEthAddress() string {
	if m != nil {
		return m.BridgeEthAddress
	}
	return ""
}

func (m *Params) GetBridgeChainId() uint64 {
	if m != nil {
		return m.BridgeChainId
	}
	return 0
}

func (m *Params) GetSignedValsetsWindow() uint64 {
	if m != nil {
		return m.SignedValsetsWindow
	}
	return 0
}

func (m *Params) GetSignedBatchesWindow() uint64 {
	if m != nil {
		return m.SignedBatchesWindow
	}
	return 0
}

func (m *Params) GetSignedClaimsWindow() uint64 {
	if m != nil {
		return m.SignedClaimsWindow
	}
	return 0
}

func (m *Params) GetTargetBatchTimeout() uint64 {
	if m != nil {
		return m.TargetBatchTimeout
	}
	return 0
}

func (m *Params) GetAverageBlockTime() uint64 {
	if m != nil {
		return m.AverageBlockTime
	}
	return 0
}

func (m *Params) GetAverageEthBlockTime() uint64 {
	if m != nil {
		return m.AverageEthBlockTime
	}
	return 0
}

func (m *Params) GetUnbondSlashingValsetsWindow() uint64 {
	if m != nil {
		return m.UnbondSlashingValsetsWindow
	}
	return 0
}

func (m *Params) GetIbcTransferTimeoutHeight() uint64 {
	if m != nil {
		return m.IbcTransferTimeoutHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "fx.gravity.v1.Params")
}

func init() { proto.RegisterFile("gravity/v1/params.proto", fileDescriptor_8772bac9489530eb) }

var fileDescriptor_8772bac9489530eb = []byte{
	// 623 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xcf, 0x4f, 0x13, 0x41,
	0x14, 0xc7, 0xbb, 0x06, 0x51, 0x46, 0x10, 0x9c, 0x16, 0xdd, 0x40, 0x58, 0x1a, 0x0e, 0x84, 0x03,
	0x74, 0x41, 0x6e, 0x26, 0x1e, 0x6c, 0x85, 0xc0, 0x8d, 0x00, 0x6a, 0xe2, 0x65, 0x9c, 0x9d, 0x9d,
	0xee, 0x4e, 0x68, 0x67, 0x9a, 0x99, 0xd9, 0xb6, 0x78, 0xf2, 0xe2, 0xdd, 0x3f, 0x8b, 0x23, 0x47,
	0x63, 0x0c, 0x31, 0xf0, 0x8f, 0x98, 0x7d, 0x33, 0x85, 0x42, 0x3c, 0x35, 0x9e, 0xda, 0xbc, 0xcf,
	0xf7, 0x47, 0xf7, 0xcd, 0x74, 0xd1, 0xab, 0x4c, 0xd3, 0xbe, 0xb0, 0xe7, 0x71, 0x7f, 0x27, 0xee,
	0x51, 0x4d, 0xbb, 0xa6, 0xd1, 0xd3, 0xca, 0x2a, 0x3c, 0xd7, 0x1e, 0x36, 0x3c, 0x6b, 0xf4, 0x77,
	0x96, 0x6a, 0x99, 0xca, 0x14, 0x90, 0xb8, 0xfc, 0xe6, 0x44, 0x4b, 0xd5, 0x31, 0xb7, 0x1d, 0xba,
	0xe1, 0xda, 0xf7, 0x19, 0x34, 0x7d, 0x04, 0x51, 0x78, 0x05, 0x21, 0xaf, 0x20, 0x22, 0x0d, 0x83,
	0x7a, 0xb0, 0x31, 0x73, 0x3c, 0xe3, 0x27, 0x87, 0x29, 0xde, 0x46, 0x35, 0xa6, 0xa4, 0xd5, 0x94,
	0x59, 0x62, 0x54, 0xa1, 0x19, 0x27, 0x39, 0x35, 0x79, 0xf8, 0x08, 0x84, 0x78, 0xc4, 0x4e, 0x00,
	0x1d, 0x50, 0x93, 0xe3, 0x4d, 0x84, 0x13, 0x2d, 0xd2, 0x8c, 0x13, 0x6e, 0x73, 0x42, 0xd3, 0x54,
	0x73, 0x63, 0xc2, 0x29, 0xd0, 0x2f, 0x38, 0xb2, 0x67, 0xf3, 0x77, 0x6e, 0x8e, 0xd7, 0xd1, 0xbc,
	0x57, 0xb3, 0x9c, 0x0a, 0x59, 0xfe, 0x86, 0xc7, 0xf5, 0x60, 0x63, 0xea, 0x78, 0xce, 0x8d, 0x5b,
	0xe5, 0xf4, 0x30, 0xc5, 0xaf, 0xd1, 0xa2, 0x11, 0x99, 0xe4, 0x29, 0xe9, 0xd3, 0x8e, 0xe1, 0xd6,
	0x90, 0x81, 0x90, 0xa9, 0x1a, 0x84, 0xd3, 0xa0, 0xae, 0x3a, 0xf8, 0xd1, 0xb1, 0x4f, 0x80, 0xc6,
	0x3c, 0x09, 0xb5, 0x2c, 0xe7, 0xb7, 0x9e, 0x27, 0xe3, 0x9e, 0xa6, 0x63, 0xde, 0xb3, 0x8d, 0x6a,
	0xde, 0xc3, 0x3a, 0x54, 0x74, 0x6f, 0x2d, 0x4f, 0xc1, 0x82, 0x1d, 0x6b, 0x01, 0xba, 0x73, 0x58,
	0xaa, 0x33, 0x6e, 0x5d, 0x0b, 0xb1, 0xa2, 0xcb, 0x55, 0x61, 0x43, 0xe4, 0x1c, 0x8e, 0x41, 0xc9,
	0xa9, 0x23, 0xe5, 0x86, 0x68, 0x9f, 0x6b, 0x9a, 0x71, 0x92, 0x74, 0x14, 0x3b, 0x03, 0x4b, 0xf8,
	0x0c, 0xf4, 0x0b, 0x9e, 0x34, 0x4b, 0x50, 0x1a, 0xf0, 0x2e, 0x7a, 0x39, 0x52, 0x97, 0x0b, 0x1d,
	0x73, 0xcc, 0xba, 0xc7, 0xf0, 0x74, 0xcf, 0xe6, 0x77, 0xa6, 0x04, 0x2d, 0x9a, 0x0e, 0x35, 0x39,
	0x69, 0x97, 0xa7, 0x23, 0x94, 0xf4, 0x6b, 0x0b, 0xe7, 0xea, 0xc1, 0xc6, 0x6c, 0xb3, 0x71, 0x71,
	0xb5, 0x5a, 0xf9, 0x75, 0xb5, 0xba, 0x9e, 0x09, 0x9b, 0x17, 0x49, 0x83, 0xa9, 0x6e, 0xcc, 0x94,
	0xe9, 0x2a, 0xe3, 0x3f, 0xb6, 0x4c, 0x7a, 0x16, 0xdb, 0xf3, 0x1e, 0x37, 0x8d, 0xf7, 0x9c, 0x1d,
	0x57, 0x21, 0x6c, 0xdf, 0x67, 0xb9, 0x2d, 0xe3, 0x2f, 0xa8, 0xf6, 0xa0, 0x03, 0x16, 0x10, 0x3e,
	0x9f, 0xa8, 0x02, 0xdf, 0xab, 0x80, 0x7d, 0xfd, 0xa3, 0x01, 0x0e, 0x25, 0x9c, 0xff, 0x0f, 0x0d,
	0x70, 0x86, 0x78, 0x80, 0xea, 0x0f, 0x1b, 0x94, 0x6c, 0x77, 0x04, 0xb3, 0x42, 0x66, 0xbe, 0x6d,
	0x61, 0xa2, 0xb6, 0x95, 0xfb, 0x6d, 0x77, 0xa9, 0xae, 0xb8, 0x85, 0xa2, 0x42, 0x26, 0x4a, 0xa6,
	0x04, 0x74, 0x65, 0xdb, 0x83, 0x8b, 0xfd, 0x02, 0x4e, 0x77, 0xd9, 0xa9, 0x4e, 0xbc, 0xe8, 0xfe,
	0x05, 0x7f, 0x8b, 0x96, 0x45, 0xc2, 0x88, 0xd5, 0x54, 0x9a, 0x36, 0xd7, 0xa3, 0xab, 0x47, 0x72,
	0x2e, 0xb2, 0xdc, 0x86, 0x18, 0x12, 0x42, 0x91, 0xb0, 0x53, 0xaf, 0xf0, 0x37, 0xf0, 0x00, 0x38,
	0xfe, 0x8a, 0xd6, 0x5c, 0x27, 0x29, 0x7a, 0x29, 0xb5, 0x9c, 0xf4, 0xd4, 0x80, 0xeb, 0xf2, 0x8f,
	0x28, 0x33, 0x4e, 0x7a, 0x5c, 0x33, 0x2e, 0x6d, 0x58, 0x9d, 0xe8, 0xf1, 0x23, 0x97, 0xfc, 0x01,
	0x82, 0x8f, 0xca, 0xdc, 0x16, 0xc4, 0x1e, 0xb9, 0xd4, 0x37, 0x53, 0xdf, 0x7e, 0xd7, 0x2b, 0xcd,
	0xfd, 0x8b, 0xeb, 0x28, 0xb8, 0xbc, 0x8e, 0x82, 0x3f, 0xd7, 0x51, 0xf0, 0xe3, 0x26, 0xaa, 0x5c,
	0xde, 0x44, 0x95, 0x9f, 0x37, 0x51, 0xe5, 0xf3, 0xe6, 0x58, 0x4f, 0xbb, 0x90, 0xb0, 0xc4, 0x61,
	0xdc, 0x1e, 0x6e, 0x31, 0xa5, 0x79, 0x3c, 0x8c, 0x47, 0x6f, 0x35, 0x68, 0x4c, 0xa6, 0xe1, 0xb5,
	0xb6, 0xfb, 0x37, 0x00, 0x00, 0xff, 0xff, 0x12, 0x05, 0x92, 0xab, 0x2b, 0x05, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ValsetUpdatePowerChangePercent.Size()
		i -= size
		if _, err := m.ValsetUpdatePowerChangePercent.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x9a
	if m.IbcTransferTimeoutHeight != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.IbcTransferTimeoutHeight))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x90
	}
	if m.UnbondSlashingValsetsWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.UnbondSlashingValsetsWindow))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x88
	}
	{
		size := m.SlashFractionConflictingClaim.Size()
		i -= size
		if _, err := m.SlashFractionConflictingClaim.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x82
	{
		size := m.SlashFractionClaim.Size()
		i -= size
		if _, err := m.SlashFractionClaim.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x7a
	{
		size := m.SlashFractionBatch.Size()
		i -= size
		if _, err := m.SlashFractionBatch.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x72
	{
		size := m.SlashFractionValset.Size()
		i -= size
		if _, err := m.SlashFractionValset.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x6a
	if m.AverageEthBlockTime != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.AverageEthBlockTime))
		i--
		dAtA[i] = 0x60
	}
	if m.AverageBlockTime != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.AverageBlockTime))
		i--
		dAtA[i] = 0x58
	}
	if m.TargetBatchTimeout != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.TargetBatchTimeout))
		i--
		dAtA[i] = 0x50
	}
	if m.SignedClaimsWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SignedClaimsWindow))
		i--
		dAtA[i] = 0x40
	}
	if m.SignedBatchesWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SignedBatchesWindow))
		i--
		dAtA[i] = 0x38
	}
	if m.SignedValsetsWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SignedValsetsWindow))
		i--
		dAtA[i] = 0x30
	}
	if m.BridgeChainId != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.BridgeChainId))
		i--
		dAtA[i] = 0x28
	}
	if len(m.BridgeEthAddress) > 0 {
		i -= len(m.BridgeEthAddress)
		copy(dAtA[i:], m.BridgeEthAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BridgeEthAddress)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ContractSourceHash) > 0 {
		i -= len(m.ContractSourceHash)
		copy(dAtA[i:], m.ContractSourceHash)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ContractSourceHash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.GravityId) > 0 {
		i -= len(m.GravityId)
		copy(dAtA[i:], m.GravityId)
		i = encodeVarintParams(dAtA, i, uint64(len(m.GravityId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GravityId)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.ContractSourceHash)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.BridgeEthAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.BridgeChainId != 0 {
		n += 1 + sovParams(uint64(m.BridgeChainId))
	}
	if m.SignedValsetsWindow != 0 {
		n += 1 + sovParams(uint64(m.SignedValsetsWindow))
	}
	if m.SignedBatchesWindow != 0 {
		n += 1 + sovParams(uint64(m.SignedBatchesWindow))
	}
	if m.SignedClaimsWindow != 0 {
		n += 1 + sovParams(uint64(m.SignedClaimsWindow))
	}
	if m.TargetBatchTimeout != 0 {
		n += 1 + sovParams(uint64(m.TargetBatchTimeout))
	}
	if m.AverageBlockTime != 0 {
		n += 1 + sovParams(uint64(m.AverageBlockTime))
	}
	if m.AverageEthBlockTime != 0 {
		n += 1 + sovParams(uint64(m.AverageEthBlockTime))
	}
	l = m.SlashFractionValset.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.SlashFractionBatch.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.SlashFractionClaim.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.SlashFractionConflictingClaim.Size()
	n += 2 + l + sovParams(uint64(l))
	if m.UnbondSlashingValsetsWindow != 0 {
		n += 2 + sovParams(uint64(m.UnbondSlashingValsetsWindow))
	}
	if m.IbcTransferTimeoutHeight != 0 {
		n += 2 + sovParams(uint64(m.IbcTransferTimeoutHeight))
	}
	l = m.ValsetUpdatePowerChangePercent.Size()
	n += 2 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GravityId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GravityId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractSourceHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractSourceHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BridgeEthAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BridgeEthAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BridgeChainId", wireType)
			}
			m.BridgeChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BridgeChainId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedValsetsWindow", wireType)
			}
			m.SignedValsetsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedValsetsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedBatchesWindow", wireType)
			}
			m.SignedBatchesWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedBatchesWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedClaimsWindow", wireType)
			}
			m.SignedClaimsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedClaimsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TargetBatchTimeout", wireType)
			}
			m.TargetBatchTimeout = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TargetBatchTimeout |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AverageBlockTime", wireType)
			}
			m.AverageBlockTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AverageBlockTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AverageEthBlockTime", wireType)
			}
			m.AverageEthBlockTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AverageEthBlockTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionValset", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionValset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionBatch", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionBatch.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionClaim", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionClaim.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionConflictingClaim", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionConflictingClaim.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 17:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondSlashingValsetsWindow", wireType)
			}
			m.UnbondSlashingValsetsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnbondSlashingValsetsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 18:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcTransferTimeoutHeight", wireType)
			}
			m.IbcTransferTimeoutHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IbcTransferTimeoutHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 19:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValsetUpdatePowerChangePercent", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ValsetUpdatePowerChangePercent.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
