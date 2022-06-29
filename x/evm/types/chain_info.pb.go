// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: evm/chain_info.proto

package types

import (
	fmt "fmt"
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

type ChainInfo struct {
	ChainID           string `protobuf:"bytes,1,opt,name=chainID,proto3" json:"chainID,omitempty"`
	SmartContractID   string `protobuf:"bytes,2,opt,name=smartContractID,proto3" json:"smartContractID,omitempty"`
	SmartContractAddr string `protobuf:"bytes,3,opt,name=smartContractAddr,proto3" json:"smartContractAddr,omitempty"`
}

func (m *ChainInfo) Reset()         { *m = ChainInfo{} }
func (m *ChainInfo) String() string { return proto.CompactTextString(m) }
func (*ChainInfo) ProtoMessage()    {}
func (*ChainInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_439fc96864abe15f, []int{0}
}
func (m *ChainInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainInfo.Merge(m, src)
}
func (m *ChainInfo) XXX_Size() int {
	return m.Size()
}
func (m *ChainInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ChainInfo proto.InternalMessageInfo

func (m *ChainInfo) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

func (m *ChainInfo) GetSmartContractID() string {
	if m != nil {
		return m.SmartContractID
	}
	return ""
}

func (m *ChainInfo) GetSmartContractAddr() string {
	if m != nil {
		return m.SmartContractAddr
	}
	return ""
}

func init() {
	proto.RegisterType((*ChainInfo)(nil), "palomachain.paloma.evm.ChainInfo")
}

func init() { proto.RegisterFile("evm/chain_info.proto", fileDescriptor_439fc96864abe15f) }

var fileDescriptor_439fc96864abe15f = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x49, 0x2d, 0xcb, 0xd5,
	0x4f, 0xce, 0x48, 0xcc, 0xcc, 0x8b, 0xcf, 0xcc, 0x4b, 0xcb, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x12, 0x2b, 0x48, 0xcc, 0xc9, 0xcf, 0x4d, 0x04, 0x8b, 0xeb, 0x41, 0xd8, 0x7a, 0xa9, 0x65,
	0xb9, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x25, 0xfa, 0x20, 0x16, 0x44, 0xb5, 0x52, 0x2d,
	0x17, 0xa7, 0x33, 0x48, 0xa5, 0x67, 0x5e, 0x5a, 0xbe, 0x90, 0x04, 0x17, 0x3b, 0x58, 0x9b, 0xa7,
	0x8b, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8c, 0x2b, 0xa4, 0xc1, 0xc5, 0x5f, 0x9c, 0x9b,
	0x58, 0x54, 0xe2, 0x9c, 0x9f, 0x57, 0x52, 0x94, 0x98, 0x5c, 0xe2, 0xe9, 0x22, 0xc1, 0x04, 0x56,
	0x81, 0x2e, 0x2c, 0xa4, 0xc3, 0x25, 0x88, 0x22, 0xe4, 0x98, 0x92, 0x52, 0x24, 0xc1, 0x0c, 0x56,
	0x8b, 0x29, 0xe1, 0xe4, 0x7c, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9,
	0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0x9a,
	0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x48, 0x3e, 0x82, 0xb2, 0xf5,
	0x2b, 0xf4, 0x41, 0x9e, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0x4e, 0x62, 0x03, 0x7b, 0xc5, 0x18, 0x10,
	0x00, 0x00, 0xff, 0xff, 0x44, 0x1b, 0xac, 0xff, 0x10, 0x01, 0x00, 0x00,
}

func (m *ChainInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SmartContractAddr) > 0 {
		i -= len(m.SmartContractAddr)
		copy(dAtA[i:], m.SmartContractAddr)
		i = encodeVarintChainInfo(dAtA, i, uint64(len(m.SmartContractAddr)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SmartContractID) > 0 {
		i -= len(m.SmartContractID)
		copy(dAtA[i:], m.SmartContractID)
		i = encodeVarintChainInfo(dAtA, i, uint64(len(m.SmartContractID)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ChainID) > 0 {
		i -= len(m.ChainID)
		copy(dAtA[i:], m.ChainID)
		i = encodeVarintChainInfo(dAtA, i, uint64(len(m.ChainID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintChainInfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovChainInfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChainInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainID)
	if l > 0 {
		n += 1 + l + sovChainInfo(uint64(l))
	}
	l = len(m.SmartContractID)
	if l > 0 {
		n += 1 + l + sovChainInfo(uint64(l))
	}
	l = len(m.SmartContractAddr)
	if l > 0 {
		n += 1 + l + sovChainInfo(uint64(l))
	}
	return n
}

func sovChainInfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozChainInfo(x uint64) (n int) {
	return sovChainInfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChainInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChainInfo
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
			return fmt.Errorf("proto: ChainInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainInfo
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
				return ErrInvalidLengthChainInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChainInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SmartContractID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainInfo
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
				return ErrInvalidLengthChainInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChainInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SmartContractID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SmartContractAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainInfo
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
				return ErrInvalidLengthChainInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChainInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SmartContractAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipChainInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthChainInfo
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
func skipChainInfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowChainInfo
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
					return 0, ErrIntOverflowChainInfo
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
					return 0, ErrIntOverflowChainInfo
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
				return 0, ErrInvalidLengthChainInfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupChainInfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthChainInfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthChainInfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowChainInfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupChainInfo = fmt.Errorf("proto: unexpected end of group")
)