// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: palomachain/paloma/skyway/bridge_transfer_limit_proposal.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
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

type SetBridgeTransferLimitProposal struct {
	Title           string                `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description     string                `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Token           string                `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	Limit           cosmossdk_io_math.Int `protobuf:"bytes,4,opt,name=limit,proto3,customtype=cosmossdk.io/math.Int" json:"limit"`
	LimitPeriod     LimitPeriod           `protobuf:"varint,5,opt,name=limit_period,json=limitPeriod,proto3,enum=palomachain.paloma.skyway.LimitPeriod" json:"limit_period,omitempty"`
	ExemptAddresses []string              `protobuf:"bytes,6,rep,name=exempt_addresses,json=exemptAddresses,proto3" json:"exempt_addresses,omitempty"`
}

func (m *SetBridgeTransferLimitProposal) Reset()         { *m = SetBridgeTransferLimitProposal{} }
func (m *SetBridgeTransferLimitProposal) String() string { return proto.CompactTextString(m) }
func (*SetBridgeTransferLimitProposal) ProtoMessage()    {}
func (*SetBridgeTransferLimitProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_e913872479061690, []int{0}
}
func (m *SetBridgeTransferLimitProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetBridgeTransferLimitProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetBridgeTransferLimitProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetBridgeTransferLimitProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetBridgeTransferLimitProposal.Merge(m, src)
}
func (m *SetBridgeTransferLimitProposal) XXX_Size() int {
	return m.Size()
}
func (m *SetBridgeTransferLimitProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SetBridgeTransferLimitProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SetBridgeTransferLimitProposal proto.InternalMessageInfo

func (m *SetBridgeTransferLimitProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SetBridgeTransferLimitProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SetBridgeTransferLimitProposal) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SetBridgeTransferLimitProposal) GetLimitPeriod() LimitPeriod {
	if m != nil {
		return m.LimitPeriod
	}
	return LimitPeriod_NONE
}

func (m *SetBridgeTransferLimitProposal) GetExemptAddresses() []string {
	if m != nil {
		return m.ExemptAddresses
	}
	return nil
}

func init() {
	proto.RegisterType((*SetBridgeTransferLimitProposal)(nil), "palomachain.paloma.skyway.SetBridgeTransferLimitProposal")
}

func init() {
	proto.RegisterFile("palomachain/paloma/skyway/bridge_transfer_limit_proposal.proto", fileDescriptor_e913872479061690)
}

var fileDescriptor_e913872479061690 = []byte{
	// 342 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0x31, 0x4b, 0xf3, 0x40,
	0x18, 0xc7, 0x93, 0xf6, 0x6d, 0xa1, 0xe9, 0x8b, 0x4a, 0xa8, 0x10, 0x0b, 0x5e, 0x83, 0x83, 0x54,
	0x90, 0x0b, 0x58, 0x5c, 0x05, 0x3b, 0x08, 0x05, 0x07, 0x89, 0x4e, 0x2e, 0xe1, 0x9a, 0x9c, 0xe9,
	0xd1, 0x24, 0xcf, 0x71, 0x77, 0x62, 0xfb, 0x2d, 0x5c, 0xfd, 0x46, 0x1d, 0x3b, 0x8a, 0x43, 0x91,
	0xf6, 0x8b, 0x48, 0xee, 0x52, 0xec, 0xa0, 0x83, 0xdb, 0xf3, 0xfc, 0xf3, 0xfb, 0xf1, 0xe4, 0xfe,
	0xce, 0x15, 0x27, 0x19, 0xe4, 0x24, 0x9e, 0x10, 0x56, 0x04, 0x66, 0x0e, 0xe4, 0x74, 0xfe, 0x42,
	0xe6, 0xc1, 0x58, 0xb0, 0x24, 0xa5, 0x91, 0x12, 0xa4, 0x90, 0x4f, 0x54, 0x44, 0x19, 0xcb, 0x99,
	0x8a, 0xb8, 0x00, 0x0e, 0x92, 0x64, 0x98, 0x0b, 0x50, 0xe0, 0x1e, 0xed, 0xf8, 0xd8, 0xcc, 0xd8,
	0xf8, 0xdd, 0x4e, 0x0a, 0x29, 0x68, 0x2a, 0x28, 0x27, 0x23, 0x74, 0x2f, 0xff, 0x78, 0xd0, 0x68,
	0x27, 0x6f, 0x35, 0x07, 0xdd, 0x53, 0x35, 0xd4, 0xc8, 0x43, 0x45, 0xdc, 0x96, 0xc0, 0x5d, 0xf5,
	0x43, 0x6e, 0xc7, 0x69, 0x28, 0xa6, 0x32, 0xea, 0xd9, 0xbe, 0xdd, 0x6f, 0x85, 0x66, 0x71, 0x7d,
	0xa7, 0x9d, 0x50, 0x19, 0x0b, 0xc6, 0x15, 0x83, 0xc2, 0xab, 0xe9, 0x6f, 0xbb, 0x91, 0xf6, 0x60,
	0x4a, 0x0b, 0xaf, 0x5e, 0x79, 0xe5, 0xe2, 0x0e, 0x9c, 0x86, 0xbe, 0xef, 0xfd, 0x2b, 0xd3, 0xe1,
	0xf1, 0x62, 0xd5, 0xb3, 0x3e, 0x56, 0xbd, 0xc3, 0x18, 0x64, 0x0e, 0x52, 0x26, 0x53, 0xcc, 0x20,
	0xc8, 0x89, 0x9a, 0xe0, 0x51, 0xa1, 0x42, 0xc3, 0xba, 0x23, 0xe7, 0x7f, 0xd5, 0x12, 0x15, 0x0c,
	0x12, 0xaf, 0xe1, 0xdb, 0xfd, 0xbd, 0x8b, 0x53, 0xfc, 0x6b, 0x49, 0xd8, 0x3c, 0x41, 0xd3, 0x61,
	0x3b, 0xfb, 0x5e, 0xdc, 0x33, 0xe7, 0x80, 0xce, 0x68, 0xce, 0x55, 0x44, 0x92, 0x44, 0x50, 0x29,
	0xa9, 0xf4, 0x9a, 0x7e, 0xbd, 0xdf, 0x0a, 0xf7, 0x4d, 0x7e, 0xbd, 0x8d, 0x87, 0x37, 0x8b, 0x35,
	0xb2, 0x97, 0x6b, 0x64, 0x7f, 0xae, 0x91, 0xfd, 0xba, 0x41, 0xd6, 0x72, 0x83, 0xac, 0xf7, 0x0d,
	0xb2, 0x1e, 0xcf, 0x53, 0xa6, 0x26, 0xcf, 0x63, 0x1c, 0x43, 0x1e, 0xfc, 0xd0, 0xfb, 0x6c, 0xdb,
	0xbc, 0x9a, 0x73, 0x2a, 0xc7, 0x4d, 0x5d, 0xf5, 0xe0, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x75, 0x71,
	0x55, 0x8e, 0x14, 0x02, 0x00, 0x00,
}

func (m *SetBridgeTransferLimitProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetBridgeTransferLimitProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetBridgeTransferLimitProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ExemptAddresses) > 0 {
		for iNdEx := len(m.ExemptAddresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExemptAddresses[iNdEx])
			copy(dAtA[i:], m.ExemptAddresses[iNdEx])
			i = encodeVarintBridgeTransferLimitProposal(dAtA, i, uint64(len(m.ExemptAddresses[iNdEx])))
			i--
			dAtA[i] = 0x32
		}
	}
	if m.LimitPeriod != 0 {
		i = encodeVarintBridgeTransferLimitProposal(dAtA, i, uint64(m.LimitPeriod))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.Limit.Size()
		i -= size
		if _, err := m.Limit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBridgeTransferLimitProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintBridgeTransferLimitProposal(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintBridgeTransferLimitProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintBridgeTransferLimitProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBridgeTransferLimitProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovBridgeTransferLimitProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SetBridgeTransferLimitProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovBridgeTransferLimitProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovBridgeTransferLimitProposal(uint64(l))
	}
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovBridgeTransferLimitProposal(uint64(l))
	}
	l = m.Limit.Size()
	n += 1 + l + sovBridgeTransferLimitProposal(uint64(l))
	if m.LimitPeriod != 0 {
		n += 1 + sovBridgeTransferLimitProposal(uint64(m.LimitPeriod))
	}
	if len(m.ExemptAddresses) > 0 {
		for _, s := range m.ExemptAddresses {
			l = len(s)
			n += 1 + l + sovBridgeTransferLimitProposal(uint64(l))
		}
	}
	return n
}

func sovBridgeTransferLimitProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBridgeTransferLimitProposal(x uint64) (n int) {
	return sovBridgeTransferLimitProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SetBridgeTransferLimitProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridgeTransferLimitProposal
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
			return fmt.Errorf("proto: SetBridgeTransferLimitProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetBridgeTransferLimitProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimitProposal
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
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimitProposal
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
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimitProposal
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
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimitProposal
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
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Limit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LimitPeriod", wireType)
			}
			m.LimitPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimitProposal
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
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExemptAddresses", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTransferLimitProposal
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
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTransferLimitProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExemptAddresses = append(m.ExemptAddresses, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBridgeTransferLimitProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBridgeTransferLimitProposal
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
func skipBridgeTransferLimitProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBridgeTransferLimitProposal
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
					return 0, ErrIntOverflowBridgeTransferLimitProposal
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
					return 0, ErrIntOverflowBridgeTransferLimitProposal
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
				return 0, ErrInvalidLengthBridgeTransferLimitProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBridgeTransferLimitProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBridgeTransferLimitProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBridgeTransferLimitProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBridgeTransferLimitProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBridgeTransferLimitProposal = fmt.Errorf("proto: unexpected end of group")
)
