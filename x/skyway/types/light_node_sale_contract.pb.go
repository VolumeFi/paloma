// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: palomachain/paloma/skyway/light_node_sale_contract.proto

package types

import (
	fmt "fmt"
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

// LightNodeSaleContract defines parameters for each of the node sale
// smart contracts deployed in external chains.
// For now, we use the address of each one to validate the message origin.
type LightNodeSaleContract struct {
	ChainReferenceId string `protobuf:"bytes,1,opt,name=chain_reference_id,json=chainReferenceId,proto3" json:"chain_reference_id,omitempty"`
	ContractAddress  string `protobuf:"bytes,2,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
}

func (m *LightNodeSaleContract) Reset()         { *m = LightNodeSaleContract{} }
func (m *LightNodeSaleContract) String() string { return proto.CompactTextString(m) }
func (*LightNodeSaleContract) ProtoMessage()    {}
func (*LightNodeSaleContract) Descriptor() ([]byte, []int) {
	return fileDescriptor_c24364f3379af24c, []int{0}
}
func (m *LightNodeSaleContract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LightNodeSaleContract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LightNodeSaleContract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LightNodeSaleContract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LightNodeSaleContract.Merge(m, src)
}
func (m *LightNodeSaleContract) XXX_Size() int {
	return m.Size()
}
func (m *LightNodeSaleContract) XXX_DiscardUnknown() {
	xxx_messageInfo_LightNodeSaleContract.DiscardUnknown(m)
}

var xxx_messageInfo_LightNodeSaleContract proto.InternalMessageInfo

func (m *LightNodeSaleContract) GetChainReferenceId() string {
	if m != nil {
		return m.ChainReferenceId
	}
	return ""
}

func (m *LightNodeSaleContract) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*LightNodeSaleContract)(nil), "palomachain.paloma.skyway.LightNodeSaleContract")
}

func init() {
	proto.RegisterFile("palomachain/paloma/skyway/light_node_sale_contract.proto", fileDescriptor_c24364f3379af24c)
}

var fileDescriptor_c24364f3379af24c = []byte{
	// 220 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0x28, 0x48, 0xcc, 0xc9,
	0xcf, 0x4d, 0x4c, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x87, 0xb0, 0xf5, 0x8b, 0xb3, 0x2b, 0xcb, 0x13,
	0x2b, 0xf5, 0x73, 0x32, 0xd3, 0x33, 0x4a, 0xe2, 0xf3, 0xf2, 0x53, 0x52, 0xe3, 0x8b, 0x13, 0x73,
	0x52, 0xe3, 0x93, 0xf3, 0xf3, 0x4a, 0x8a, 0x12, 0x93, 0x4b, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2,
	0x85, 0x24, 0x91, 0x74, 0xea, 0x41, 0xd8, 0x7a, 0x10, 0x9d, 0x4a, 0x05, 0x5c, 0xa2, 0x3e, 0x20,
	0xcd, 0x7e, 0xf9, 0x29, 0xa9, 0xc1, 0x89, 0x39, 0xa9, 0xce, 0x50, 0x9d, 0x42, 0x3a, 0x5c, 0x42,
	0x60, 0xf5, 0xf1, 0x45, 0xa9, 0x69, 0xa9, 0x45, 0xa9, 0x79, 0xc9, 0xa9, 0xf1, 0x99, 0x29, 0x12,
	0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x02, 0x60, 0x99, 0x20, 0x98, 0x84, 0x67, 0x8a, 0x90, 0x26,
	0x97, 0x00, 0xcc, 0xce, 0xf8, 0xc4, 0x94, 0x94, 0xa2, 0xd4, 0xe2, 0x62, 0x09, 0x26, 0xb0, 0x5a,
	0x7e, 0x98, 0xb8, 0x23, 0x44, 0xd8, 0xc9, 0xed, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18,
	0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5,
	0x18, 0xa2, 0x74, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xb1, 0xf8,
	0xb5, 0x02, 0xe6, 0xdb, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0xb0, 0xdf, 0x8c, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xc9, 0x39, 0xe4, 0xa9, 0x17, 0x01, 0x00, 0x00,
}

func (m *LightNodeSaleContract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LightNodeSaleContract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LightNodeSaleContract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ContractAddress) > 0 {
		i -= len(m.ContractAddress)
		copy(dAtA[i:], m.ContractAddress)
		i = encodeVarintLightNodeSaleContract(dAtA, i, uint64(len(m.ContractAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ChainReferenceId) > 0 {
		i -= len(m.ChainReferenceId)
		copy(dAtA[i:], m.ChainReferenceId)
		i = encodeVarintLightNodeSaleContract(dAtA, i, uint64(len(m.ChainReferenceId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintLightNodeSaleContract(dAtA []byte, offset int, v uint64) int {
	offset -= sovLightNodeSaleContract(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LightNodeSaleContract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainReferenceId)
	if l > 0 {
		n += 1 + l + sovLightNodeSaleContract(uint64(l))
	}
	l = len(m.ContractAddress)
	if l > 0 {
		n += 1 + l + sovLightNodeSaleContract(uint64(l))
	}
	return n
}

func sovLightNodeSaleContract(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLightNodeSaleContract(x uint64) (n int) {
	return sovLightNodeSaleContract(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LightNodeSaleContract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLightNodeSaleContract
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
			return fmt.Errorf("proto: LightNodeSaleContract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LightNodeSaleContract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainReferenceId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLightNodeSaleContract
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
				return ErrInvalidLengthLightNodeSaleContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLightNodeSaleContract
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainReferenceId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLightNodeSaleContract
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
				return ErrInvalidLengthLightNodeSaleContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLightNodeSaleContract
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLightNodeSaleContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLightNodeSaleContract
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
func skipLightNodeSaleContract(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLightNodeSaleContract
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
					return 0, ErrIntOverflowLightNodeSaleContract
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
					return 0, ErrIntOverflowLightNodeSaleContract
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
				return 0, ErrInvalidLengthLightNodeSaleContract
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLightNodeSaleContract
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLightNodeSaleContract
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLightNodeSaleContract        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLightNodeSaleContract          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLightNodeSaleContract = fmt.Errorf("proto: unexpected end of group")
)