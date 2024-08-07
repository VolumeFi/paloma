// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: palomachain/paloma/skyway/bridge_tax_proposal.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

type SetBridgeTaxProposal struct {
	Title           string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description     string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Rate            string   `protobuf:"bytes,3,opt,name=rate,proto3" json:"rate,omitempty"`
	ExemptAddresses []string `protobuf:"bytes,5,rep,name=exempt_addresses,json=exemptAddresses,proto3" json:"exempt_addresses,omitempty"`
	Token           string   `protobuf:"bytes,6,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *SetBridgeTaxProposal) Reset()         { *m = SetBridgeTaxProposal{} }
func (m *SetBridgeTaxProposal) String() string { return proto.CompactTextString(m) }
func (*SetBridgeTaxProposal) ProtoMessage()    {}
func (*SetBridgeTaxProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d240d8fd882781b, []int{0}
}
func (m *SetBridgeTaxProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetBridgeTaxProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetBridgeTaxProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetBridgeTaxProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetBridgeTaxProposal.Merge(m, src)
}
func (m *SetBridgeTaxProposal) XXX_Size() int {
	return m.Size()
}
func (m *SetBridgeTaxProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SetBridgeTaxProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SetBridgeTaxProposal proto.InternalMessageInfo

func (m *SetBridgeTaxProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SetBridgeTaxProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SetBridgeTaxProposal) GetRate() string {
	if m != nil {
		return m.Rate
	}
	return ""
}

func (m *SetBridgeTaxProposal) GetExemptAddresses() []string {
	if m != nil {
		return m.ExemptAddresses
	}
	return nil
}

func (m *SetBridgeTaxProposal) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*SetBridgeTaxProposal)(nil), "palomachain.paloma.skyway.SetBridgeTaxProposal")
}

func init() {
	proto.RegisterFile("palomachain/paloma/skyway/bridge_tax_proposal.proto", fileDescriptor_2d240d8fd882781b)
}

var fileDescriptor_2d240d8fd882781b = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x50, 0xbd, 0x4e, 0xf3, 0x30,
	0x14, 0xad, 0xbf, 0xfe, 0xe8, 0xab, 0x91, 0xa0, 0x8a, 0x3a, 0xa4, 0x0c, 0x56, 0xd5, 0xa9, 0x48,
	0x90, 0x0c, 0x7d, 0x02, 0x2a, 0xc4, 0xc0, 0x84, 0x0a, 0x13, 0x4b, 0xe4, 0xda, 0x57, 0xad, 0xd5,
	0x24, 0xb6, 0x6c, 0x57, 0x24, 0x6f, 0xc1, 0xc3, 0xf0, 0x10, 0x88, 0xa9, 0x62, 0x62, 0x44, 0xc9,
	0x8b, 0xa0, 0xd8, 0xa9, 0xd4, 0x81, 0xed, 0x9e, 0x73, 0xee, 0xb9, 0xf7, 0xe8, 0xe0, 0x85, 0xa2,
	0xa9, 0xcc, 0x28, 0xdb, 0x52, 0x91, 0xc7, 0x7e, 0x8e, 0xcd, 0xae, 0x7c, 0xa5, 0x65, 0xbc, 0xd6,
	0x82, 0x6f, 0x20, 0xb1, 0xb4, 0x48, 0x94, 0x96, 0x4a, 0x1a, 0x9a, 0x46, 0x4a, 0x4b, 0x2b, 0x83,
	0xc9, 0x89, 0x29, 0xf2, 0x73, 0xe4, 0x4d, 0x97, 0x13, 0x26, 0x4d, 0x26, 0x4d, 0xe2, 0x16, 0x63,
	0x0f, 0xbc, 0x6b, 0xf6, 0x89, 0xf0, 0xf8, 0x09, 0xec, 0xd2, 0x9d, 0x7d, 0xa6, 0xc5, 0x63, 0x7b,
	0x34, 0x18, 0xe3, 0xbe, 0x15, 0x36, 0x85, 0x10, 0x4d, 0xd1, 0x7c, 0xb8, 0xf2, 0x20, 0x98, 0xe2,
	0x33, 0x0e, 0x86, 0x69, 0xa1, 0xac, 0x90, 0x79, 0xf8, 0xcf, 0x69, 0xa7, 0x54, 0x30, 0xc3, 0x3d,
	0x4d, 0x2d, 0x84, 0xdd, 0x46, 0x5a, 0x9e, 0x7f, 0xbd, 0xdf, 0xe0, 0xf6, 0xe1, 0x1d, 0xb0, 0x95,
	0xd3, 0x82, 0x2b, 0x3c, 0x82, 0x02, 0x32, 0x65, 0x13, 0xca, 0xb9, 0x06, 0x63, 0xc0, 0x84, 0xfd,
	0x69, 0x77, 0x3e, 0x5c, 0x5d, 0x78, 0xfe, 0xf6, 0x48, 0xbb, 0x18, 0x72, 0x07, 0x79, 0x38, 0x68,
	0x63, 0x34, 0xe0, 0xa1, 0xf7, 0xbf, 0x37, 0xea, 0x37, 0xcb, 0x2c, 0xdd, 0x73, 0xe0, 0x89, 0x63,
	0xcd, 0xf2, 0xfe, 0xa3, 0x22, 0xe8, 0x50, 0x11, 0xf4, 0x53, 0x11, 0xf4, 0x56, 0x93, 0xce, 0xa1,
	0x26, 0x9d, 0xef, 0x9a, 0x74, 0x5e, 0xae, 0x37, 0xc2, 0x6e, 0xf7, 0xeb, 0x88, 0xc9, 0x2c, 0xfe,
	0xa3, 0xdc, 0xe2, 0x58, 0xaf, 0x2d, 0x15, 0x98, 0xf5, 0xc0, 0x75, 0xb3, 0xf8, 0x0d, 0x00, 0x00,
	0xff, 0xff, 0x64, 0xd1, 0x17, 0xa7, 0x88, 0x01, 0x00, 0x00,
}

func (m *SetBridgeTaxProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetBridgeTaxProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetBridgeTaxProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintBridgeTaxProposal(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.ExemptAddresses) > 0 {
		for iNdEx := len(m.ExemptAddresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExemptAddresses[iNdEx])
			copy(dAtA[i:], m.ExemptAddresses[iNdEx])
			i = encodeVarintBridgeTaxProposal(dAtA, i, uint64(len(m.ExemptAddresses[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Rate) > 0 {
		i -= len(m.Rate)
		copy(dAtA[i:], m.Rate)
		i = encodeVarintBridgeTaxProposal(dAtA, i, uint64(len(m.Rate)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintBridgeTaxProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintBridgeTaxProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBridgeTaxProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovBridgeTaxProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SetBridgeTaxProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovBridgeTaxProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovBridgeTaxProposal(uint64(l))
	}
	l = len(m.Rate)
	if l > 0 {
		n += 1 + l + sovBridgeTaxProposal(uint64(l))
	}
	if len(m.ExemptAddresses) > 0 {
		for _, s := range m.ExemptAddresses {
			l = len(s)
			n += 1 + l + sovBridgeTaxProposal(uint64(l))
		}
	}
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovBridgeTaxProposal(uint64(l))
	}
	return n
}

func sovBridgeTaxProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBridgeTaxProposal(x uint64) (n int) {
	return sovBridgeTaxProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SetBridgeTaxProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridgeTaxProposal
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
			return fmt.Errorf("proto: SetBridgeTaxProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetBridgeTaxProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTaxProposal
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
				return ErrInvalidLengthBridgeTaxProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTaxProposal
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
					return ErrIntOverflowBridgeTaxProposal
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
				return ErrInvalidLengthBridgeTaxProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTaxProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTaxProposal
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
				return ErrInvalidLengthBridgeTaxProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTaxProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExemptAddresses", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTaxProposal
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
				return ErrInvalidLengthBridgeTaxProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTaxProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExemptAddresses = append(m.ExemptAddresses, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridgeTaxProposal
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
				return ErrInvalidLengthBridgeTaxProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridgeTaxProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBridgeTaxProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBridgeTaxProposal
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
func skipBridgeTaxProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBridgeTaxProposal
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
					return 0, ErrIntOverflowBridgeTaxProposal
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
					return 0, ErrIntOverflowBridgeTaxProposal
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
				return 0, ErrInvalidLengthBridgeTaxProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBridgeTaxProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBridgeTaxProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBridgeTaxProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBridgeTaxProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBridgeTaxProposal = fmt.Errorf("proto: unexpected end of group")
)
