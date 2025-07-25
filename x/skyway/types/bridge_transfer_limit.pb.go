// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: palomachain/paloma/skyway/bridge_transfer_limit.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type LimitPeriod int32

const (
	LimitPeriod_NONE    LimitPeriod = 0
	LimitPeriod_DAILY   LimitPeriod = 1
	LimitPeriod_WEEKLY  LimitPeriod = 2
	LimitPeriod_MONTHLY LimitPeriod = 3
	LimitPeriod_YEARLY  LimitPeriod = 4
)

var LimitPeriod_name = map[int32]string{
	0: "NONE",
	1: "DAILY",
	2: "WEEKLY",
	3: "MONTHLY",
	4: "YEARLY",
}

var LimitPeriod_value = map[string]int32{
	"NONE":    0,
	"DAILY":   1,
	"WEEKLY":  2,
	"MONTHLY": 3,
	"YEARLY":  4,
}

func (x LimitPeriod) String() string {
	return proto.EnumName(LimitPeriod_name, int32(x))
}

func (LimitPeriod) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_30c19b73abd7126a, []int{0}
}

// Allow at most `limit` tokens of `token` to be transferred within a
// `limit_period` window. `limit_period` will be converted to blocks.
// If more than `limit` tokens are attempted to be transferred between those
// block heights, the transfer is not allowed.
// If the sender is in `exempt_addresses`, the limits are not checked nor
// updated.
type BridgeTransferLimit struct {
	Token           string                                          `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Limit           cosmossdk_io_math.Int                           `protobuf:"bytes,2,opt,name=limit,proto3,customtype=cosmossdk.io/math.Int" json:"limit"`
	LimitPeriod     LimitPeriod                                     `protobuf:"varint,3,opt,name=limit_period,json=limitPeriod,proto3,enum=palomachain.paloma.skyway.LimitPeriod" json:"limit_period,omitempty"`
	ExemptAddresses []github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,4,rep,name=exempt_addresses,json=exemptAddresses,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"exempt_addresses,omitempty"`
}

func (m *BridgeTransferLimit) Reset()         { *m = BridgeTransferLimit{} }
func (m *BridgeTransferLimit) String() string { return proto.CompactTextString(m) }
func (*BridgeTransferLimit) ProtoMessage()    {}
func (*BridgeTransferLimit) Descriptor() ([]byte, []int) {
	return fileDescriptor_30c19b73abd7126a, []int{0}
}
func (m *BridgeTransferLimit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BridgeTransferLimit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BridgeTransferLimit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BridgeTransferLimit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BridgeTransferLimit.Merge(m, src)
}
func (m *BridgeTransferLimit) XXX_Size() int {
	return m.Size()
}
func (m *BridgeTransferLimit) XXX_DiscardUnknown() {
	xxx_messageInfo_BridgeTransferLimit.DiscardUnknown(m)
}

var xxx_messageInfo_BridgeTransferLimit proto.InternalMessageInfo

func (m *BridgeTransferLimit) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *BridgeTransferLimit) GetLimitPeriod() LimitPeriod {
	if m != nil {
		return m.LimitPeriod
	}
	return LimitPeriod_NONE
}

func (m *BridgeTransferLimit) GetExemptAddresses() []github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.ExemptAddresses
	}
	return nil
}

// Transfer usage counters used to check for transfer limits for a single denom.
// `total` maintains the total amount transferred during the current window.
// `start_block_height` maintains the block height of the first transfer in the
// current window.
type BridgeTransferUsage struct {
	Total            cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=total,proto3,customtype=cosmossdk.io/math.Int" json:"total"`
	StartBlockHeight int64                 `protobuf:"varint,2,opt,name=start_block_height,json=startBlockHeight,proto3" json:"start_block_height,omitempty"`
}

func (m *BridgeTransferUsage) Reset()         { *m = BridgeTransferUsage{} }
func (m *BridgeTransferUsage) String() string { return proto.CompactTextString(m) }
func (*BridgeTransferUsage) ProtoMessage()    {}
func (*BridgeTransferUsage) Descriptor() ([]byte, []int) {
	return fileDescriptor_30c19b73abd7126a, []int{1}
}
func (m *BridgeTransferUsage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BridgeTransferUsage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BridgeTransferUsage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BridgeTransferUsage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BridgeTransferUsage.Merge(m, src)
}
func (m *BridgeTransferUsage) XXX_Size() int {
	return m.Size()
}
func (m *BridgeTransferUsage) XXX_DiscardUnknown() {
	xxx_messageInfo_BridgeTransferUsage.DiscardUnknown(m)
}

var xxx_messageInfo_BridgeTransferUsage proto.InternalMessageInfo

func (m *BridgeTransferUsage) GetStartBlockHeight() int64 {
	if m != nil {
		return m.StartBlockHeight
	}
	return 0
}

func init() {
	proto.RegisterEnum("palomachain.paloma.skyway.LimitPeriod", LimitPeriod_name, LimitPeriod_value)
	proto.RegisterType((*BridgeTransferLimit)(nil), "palomachain.paloma.skyway.BridgeTransferLimit")
	proto.RegisterType((*BridgeTransferUsage)(nil), "palomachain.paloma.skyway.BridgeTransferUsage")
}

func init() {
	proto.RegisterFile("palomachain/paloma/skyway/bridge_transfer_limit.proto", fileDescriptor_30c19b73abd7126a)
}

var fileDescriptor_30c19b73abd7126a = []byte{
	// 438 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x5d, 0x8b, 0xd3, 0x40,
	0x14, 0x4d, 0xfa, 0xb1, 0xba, 0xd3, 0x45, 0xc3, 0xb8, 0x42, 0x15, 0x4c, 0xcb, 0x3e, 0x48, 0x11,
	0x77, 0x06, 0x77, 0xf1, 0x07, 0x24, 0x58, 0xdc, 0x62, 0xec, 0x4a, 0x58, 0x91, 0x88, 0x10, 0x26,
	0xc9, 0x98, 0x0c, 0xf9, 0x98, 0x90, 0x19, 0xb5, 0xfd, 0x09, 0xbe, 0xf9, 0xb3, 0xf6, 0x71, 0x1f,
	0xc5, 0x87, 0x22, 0xed, 0xbf, 0xf0, 0x49, 0x32, 0xb3, 0xc5, 0x22, 0xca, 0x3e, 0xe5, 0xcc, 0xbd,
	0xe7, 0x9c, 0xdc, 0x7b, 0xb8, 0xe0, 0x79, 0x4d, 0x0a, 0x5e, 0x92, 0x38, 0x23, 0xac, 0xc2, 0x1a,
	0x63, 0x91, 0x2f, 0xbf, 0x90, 0x25, 0x8e, 0x1a, 0x96, 0xa4, 0x34, 0x94, 0x0d, 0xa9, 0xc4, 0x47,
	0xda, 0x84, 0x05, 0x2b, 0x99, 0x44, 0x75, 0xc3, 0x25, 0x87, 0x0f, 0x76, 0x64, 0x48, 0x63, 0xa4,
	0x65, 0x0f, 0x0f, 0x53, 0x9e, 0x72, 0xc5, 0xc2, 0x2d, 0xd2, 0x82, 0xa3, 0xaf, 0x1d, 0x70, 0xcf,
	0x55, 0x86, 0x17, 0xd7, 0x7e, 0x5e, 0x6b, 0x07, 0x0f, 0x41, 0x5f, 0xf2, 0x9c, 0x56, 0x43, 0x73,
	0x6c, 0x4e, 0xf6, 0x7d, 0xfd, 0x80, 0xa7, 0xa0, 0xaf, 0xfe, 0x36, 0xec, 0xb4, 0x55, 0xf7, 0xd1,
	0xe5, 0x6a, 0x64, 0xfc, 0x58, 0x8d, 0xee, 0xc7, 0x5c, 0x94, 0x5c, 0x88, 0x24, 0x47, 0x8c, 0xe3,
	0x92, 0xc8, 0x0c, 0xcd, 0x2a, 0xe9, 0x6b, 0x2e, 0x9c, 0x81, 0x03, 0x05, 0xc2, 0x9a, 0x36, 0x8c,
	0x27, 0xc3, 0xee, 0xd8, 0x9c, 0xdc, 0x39, 0x79, 0x8c, 0xfe, 0x3b, 0x2a, 0x52, 0x23, 0xbc, 0x51,
	0x6c, 0x7f, 0x50, 0xfc, 0x79, 0xc0, 0x0f, 0xc0, 0xa2, 0x0b, 0x5a, 0xd6, 0x32, 0x24, 0x49, 0xd2,
	0x50, 0x21, 0xa8, 0x18, 0xf6, 0xc6, 0xdd, 0xc9, 0x81, 0xfb, 0xec, 0xd7, 0x6a, 0x74, 0x9c, 0x32,
	0x99, 0x7d, 0x8a, 0x50, 0xcc, 0x4b, 0xac, 0x27, 0xba, 0xfe, 0x1c, 0x8b, 0x24, 0xc7, 0x72, 0x59,
	0x53, 0x81, 0x9c, 0x38, 0x76, 0xb4, 0xd4, 0xbf, 0xab, 0xad, 0x9c, 0xad, 0xd3, 0xd1, 0xe2, 0xef,
	0x28, 0xde, 0x0a, 0x92, 0xd2, 0x76, 0x69, 0xc9, 0x25, 0x29, 0x74, 0x14, 0x37, 0x2e, 0xad, 0xb8,
	0xf0, 0x29, 0x80, 0x42, 0x92, 0x46, 0x86, 0x51, 0xc1, 0xe3, 0x3c, 0xcc, 0x28, 0x4b, 0x33, 0x1d,
	0x5b, 0xd7, 0xb7, 0x54, 0xc7, 0x6d, 0x1b, 0x67, 0xaa, 0xfe, 0xe4, 0x25, 0x18, 0xec, 0xec, 0x0c,
	0x6f, 0x83, 0xde, 0xfc, 0x7c, 0x3e, 0xb5, 0x0c, 0xb8, 0x0f, 0xfa, 0x2f, 0x9c, 0x99, 0x17, 0x58,
	0x26, 0x04, 0x60, 0xef, 0xdd, 0x74, 0xfa, 0xca, 0x0b, 0xac, 0x0e, 0x1c, 0x80, 0x5b, 0xaf, 0xcf,
	0xe7, 0x17, 0x67, 0x5e, 0x60, 0x75, 0xdb, 0x46, 0x30, 0x75, 0x7c, 0x2f, 0xb0, 0x7a, 0xee, 0xec,
	0x72, 0x6d, 0x9b, 0x57, 0x6b, 0xdb, 0xfc, 0xb9, 0xb6, 0xcd, 0x6f, 0x1b, 0xdb, 0xb8, 0xda, 0xd8,
	0xc6, 0xf7, 0x8d, 0x6d, 0xbc, 0xc7, 0x3b, 0xe1, 0xfc, 0xe3, 0xb6, 0x3e, 0x9f, 0xe0, 0xc5, 0xf6,
	0xc0, 0x54, 0x52, 0xd1, 0x9e, 0x3a, 0x90, 0xd3, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xb9,
	0xd0, 0xe8, 0x8a, 0x02, 0x00, 0x00,
}

func (m *BridgeTransferLimit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BridgeTransferLimit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BridgeTransferLimit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ExemptAddresses) > 0 {
		for iNdEx := len(m.ExemptAddresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExemptAddresses[iNdEx])
			copy(dAtA[i:], m.ExemptAddresses[iNdEx])
			i = encodeVarintBridgeTransferLimit(dAtA, i, uint64(len(m.ExemptAddresses[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if m.LimitPeriod != 0 {
		i = encodeVarintBridgeTransferLimit(dAtA, i, uint64(m.LimitPeriod))
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.Limit.Size()
		i -= size
		if _, err := m.Limit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBridgeTransferLimit(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintBridgeTransferLimit(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BridgeTransferUsage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BridgeTransferUsage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BridgeTransferUsage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.StartBlockHeight != 0 {
		i = encodeVarintBridgeTransferLimit(dAtA, i, uint64(m.StartBlockHeight))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.Total.Size()
		i -= size
		if _, err := m.Total.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBridgeTransferLimit(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintBridgeTransferLimit(dAtA []byte, offset int, v uint64) int {
	offset -= sovBridgeTransferLimit(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BridgeTransferLimit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovBridgeTransferLimit(uint64(l))
	}
	l = m.Limit.Size()
	n += 1 + l + sovBridgeTransferLimit(uint64(l))
	if m.LimitPeriod != 0 {
		n += 1 + sovBridgeTransferLimit(uint64(m.LimitPeriod))
	}
	if len(m.ExemptAddresses) > 0 {
		for _, b := range m.ExemptAddresses {
			l = len(b)
			n += 1 + l + sovBridgeTransferLimit(uint64(l))
		}
	}
	return n
}

func (m *BridgeTransferUsage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Total.Size()
	n += 1 + l + sovBridgeTransferLimit(uint64(l))
	if m.StartBlockHeight != 0 {
		n += 1 + sovBridgeTransferLimit(uint64(m.StartBlockHeight))
	}
	return n
}

func sovBridgeTransferLimit(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBridgeTransferLimit(x uint64) (n int) {
	return sovBridgeTransferLimit(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BridgeTransferLimit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridgeTransferLimit
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
			return fmt.Errorf("proto: BridgeTransferLimit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BridgeTransferLimit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimit
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
				return ErrInvalidLengthBridgeTransferLimit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimit
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
				return ErrInvalidLengthBridgeTransferLimit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Limit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LimitPeriod", wireType)
			}
			m.LimitPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimit
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LimitPeriod |= LimitPeriod(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExemptAddresses", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimit
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
				return ErrInvalidLengthBridgeTransferLimit
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExemptAddresses = append(m.ExemptAddresses, make([]byte, postIndex-iNdEx))
			copy(m.ExemptAddresses[len(m.ExemptAddresses)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBridgeTransferLimit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBridgeTransferLimit
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
func (m *BridgeTransferUsage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridgeTransferLimit
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
			return fmt.Errorf("proto: BridgeTransferUsage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BridgeTransferUsage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Total", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimit
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
				return ErrInvalidLengthBridgeTransferLimit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Total.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartBlockHeight", wireType)
			}
			m.StartBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimit
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartBlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBridgeTransferLimit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBridgeTransferLimit
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
func skipBridgeTransferLimit(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBridgeTransferLimit
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
					return 0, ErrIntOverflowBridgeTransferLimit
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
					return 0, ErrIntOverflowBridgeTransferLimit
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
				return 0, ErrInvalidLengthBridgeTransferLimit
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBridgeTransferLimit
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBridgeTransferLimit
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBridgeTransferLimit        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBridgeTransferLimit          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBridgeTransferLimit = fmt.Errorf("proto: unexpected end of group")
)
