// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: palomachain/paloma/evm/set_smart_contract_deployers_proposal.proto

package types

import (
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

type SetSmartContractDeployersProposal struct {
	Title     string                                       `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Summary   string                                       `protobuf:"bytes,2,opt,name=summary,proto3" json:"summary,omitempty"`
	Deployers []SetSmartContractDeployersProposal_Deployer `protobuf:"bytes,3,rep,name=deployers,proto3" json:"deployers"`
}

func (m *SetSmartContractDeployersProposal) Reset()         { *m = SetSmartContractDeployersProposal{} }
func (m *SetSmartContractDeployersProposal) String() string { return proto.CompactTextString(m) }
func (*SetSmartContractDeployersProposal) ProtoMessage()    {}
func (*SetSmartContractDeployersProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_98112a37abfe69b8, []int{0}
}
func (m *SetSmartContractDeployersProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetSmartContractDeployersProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetSmartContractDeployersProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetSmartContractDeployersProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetSmartContractDeployersProposal.Merge(m, src)
}
func (m *SetSmartContractDeployersProposal) XXX_Size() int {
	return m.Size()
}
func (m *SetSmartContractDeployersProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SetSmartContractDeployersProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SetSmartContractDeployersProposal proto.InternalMessageInfo

func (m *SetSmartContractDeployersProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SetSmartContractDeployersProposal) GetSummary() string {
	if m != nil {
		return m.Summary
	}
	return ""
}

func (m *SetSmartContractDeployersProposal) GetDeployers() []SetSmartContractDeployersProposal_Deployer {
	if m != nil {
		return m.Deployers
	}
	return nil
}

type SetSmartContractDeployersProposal_Deployer struct {
	ChainReferenceID string `protobuf:"bytes,1,opt,name=chainReferenceID,proto3" json:"chainReferenceID,omitempty"`
	ContractAddress  string `protobuf:"bytes,2,opt,name=contractAddress,proto3" json:"contractAddress,omitempty"`
}

func (m *SetSmartContractDeployersProposal_Deployer) Reset() {
	*m = SetSmartContractDeployersProposal_Deployer{}
}
func (m *SetSmartContractDeployersProposal_Deployer) String() string {
	return proto.CompactTextString(m)
}
func (*SetSmartContractDeployersProposal_Deployer) ProtoMessage() {}
func (*SetSmartContractDeployersProposal_Deployer) Descriptor() ([]byte, []int) {
	return fileDescriptor_98112a37abfe69b8, []int{0, 0}
}
func (m *SetSmartContractDeployersProposal_Deployer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetSmartContractDeployersProposal_Deployer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetSmartContractDeployersProposal_Deployer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetSmartContractDeployersProposal_Deployer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetSmartContractDeployersProposal_Deployer.Merge(m, src)
}
func (m *SetSmartContractDeployersProposal_Deployer) XXX_Size() int {
	return m.Size()
}
func (m *SetSmartContractDeployersProposal_Deployer) XXX_DiscardUnknown() {
	xxx_messageInfo_SetSmartContractDeployersProposal_Deployer.DiscardUnknown(m)
}

var xxx_messageInfo_SetSmartContractDeployersProposal_Deployer proto.InternalMessageInfo

func (m *SetSmartContractDeployersProposal_Deployer) GetChainReferenceID() string {
	if m != nil {
		return m.ChainReferenceID
	}
	return ""
}

func (m *SetSmartContractDeployersProposal_Deployer) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*SetSmartContractDeployersProposal)(nil), "palomachain.paloma.evm.SetSmartContractDeployersProposal")
	proto.RegisterType((*SetSmartContractDeployersProposal_Deployer)(nil), "palomachain.paloma.evm.SetSmartContractDeployersProposal.Deployer")
}

func init() {
	proto.RegisterFile("palomachain/paloma/evm/set_smart_contract_deployers_proposal.proto", fileDescriptor_98112a37abfe69b8)
}

var fileDescriptor_98112a37abfe69b8 = []byte{
	// 301 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0x93, 0x96, 0x7f, 0x35, 0x03, 0xc8, 0xaa, 0x50, 0xd4, 0xc1, 0x14, 0xa6, 0xc2, 0xe0,
	0x48, 0xf0, 0x04, 0xa4, 0x5d, 0xd8, 0x50, 0xba, 0xb1, 0x04, 0x37, 0xb9, 0xa6, 0x91, 0xe2, 0xd8,
	0xb2, 0xdd, 0x8a, 0xbc, 0x05, 0x2f, 0xc0, 0xfb, 0x74, 0xec, 0xc8, 0x84, 0x50, 0xf2, 0x22, 0xa8,
	0xf9, 0x03, 0x08, 0x2a, 0xb1, 0x7d, 0xf7, 0xf9, 0xf4, 0xf3, 0x77, 0x77, 0xc8, 0x93, 0x2c, 0x15,
	0x9c, 0x85, 0x0b, 0x96, 0x64, 0x6e, 0xad, 0x5d, 0x58, 0x71, 0x57, 0x83, 0x09, 0x34, 0x67, 0xca,
	0x04, 0xa1, 0xc8, 0x8c, 0x62, 0xa1, 0x09, 0x22, 0x90, 0xa9, 0xc8, 0x41, 0xe9, 0x40, 0x2a, 0x21,
	0x85, 0x66, 0x29, 0x95, 0x4a, 0x18, 0x81, 0xcf, 0x7e, 0x30, 0x68, 0xad, 0x29, 0xac, 0xf8, 0xa0,
	0x1f, 0x8b, 0x58, 0x54, 0x2d, 0xee, 0x56, 0xd5, 0xdd, 0x97, 0xaf, 0x1d, 0x74, 0x31, 0x05, 0x33,
	0xdd, 0xc2, 0xc7, 0x0d, 0x7b, 0xd2, 0xa2, 0x1f, 0x1a, 0x32, 0xee, 0xa3, 0x7d, 0x93, 0x98, 0x14,
	0x1c, 0x7b, 0x68, 0x8f, 0x7a, 0x7e, 0x5d, 0x60, 0x07, 0x1d, 0xea, 0x25, 0xe7, 0x4c, 0xe5, 0x4e,
	0xa7, 0xf2, 0xdb, 0x12, 0xcf, 0x51, 0xef, 0x2b, 0x9f, 0xd3, 0x1d, 0x76, 0x47, 0xc7, 0x37, 0x1e,
	0xdd, 0x9d, 0x8b, 0xfe, 0xfb, 0x3b, 0x6d, 0x1d, 0x6f, 0x6f, 0xfd, 0x7e, 0x6e, 0xf9, 0xdf, 0xe8,
	0xc1, 0x13, 0x3a, 0x6a, 0x1f, 0xf1, 0x35, 0x3a, 0xad, 0xd8, 0x3e, 0xcc, 0x41, 0x41, 0x16, 0xc2,
	0xfd, 0xa4, 0x89, 0xfb, 0xc7, 0xc7, 0x23, 0x74, 0xd2, 0x2e, 0xf2, 0x2e, 0x8a, 0x14, 0x68, 0xdd,
	0x4c, 0xf0, 0xdb, 0xf6, 0xc6, 0xeb, 0x82, 0xd8, 0x9b, 0x82, 0xd8, 0x1f, 0x05, 0xb1, 0x5f, 0x4a,
	0x62, 0x6d, 0x4a, 0x62, 0xbd, 0x95, 0xc4, 0x7a, 0xbc, 0x8a, 0x13, 0xb3, 0x58, 0xce, 0x68, 0x28,
	0xb8, 0xbb, 0xe3, 0x6c, 0xcf, 0xd5, 0xe1, 0x4c, 0x2e, 0x41, 0xcf, 0x0e, 0xaa, 0x5d, 0xdf, 0x7e,
	0x06, 0x00, 0x00, 0xff, 0xff, 0xf7, 0xe6, 0x62, 0xea, 0xdf, 0x01, 0x00, 0x00,
}

func (m *SetSmartContractDeployersProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetSmartContractDeployersProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetSmartContractDeployersProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Deployers) > 0 {
		for iNdEx := len(m.Deployers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Deployers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSetSmartContractDeployersProposal(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Summary) > 0 {
		i -= len(m.Summary)
		copy(dAtA[i:], m.Summary)
		i = encodeVarintSetSmartContractDeployersProposal(dAtA, i, uint64(len(m.Summary)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintSetSmartContractDeployersProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SetSmartContractDeployersProposal_Deployer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetSmartContractDeployersProposal_Deployer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetSmartContractDeployersProposal_Deployer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ContractAddress) > 0 {
		i -= len(m.ContractAddress)
		copy(dAtA[i:], m.ContractAddress)
		i = encodeVarintSetSmartContractDeployersProposal(dAtA, i, uint64(len(m.ContractAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ChainReferenceID) > 0 {
		i -= len(m.ChainReferenceID)
		copy(dAtA[i:], m.ChainReferenceID)
		i = encodeVarintSetSmartContractDeployersProposal(dAtA, i, uint64(len(m.ChainReferenceID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSetSmartContractDeployersProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovSetSmartContractDeployersProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SetSmartContractDeployersProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovSetSmartContractDeployersProposal(uint64(l))
	}
	l = len(m.Summary)
	if l > 0 {
		n += 1 + l + sovSetSmartContractDeployersProposal(uint64(l))
	}
	if len(m.Deployers) > 0 {
		for _, e := range m.Deployers {
			l = e.Size()
			n += 1 + l + sovSetSmartContractDeployersProposal(uint64(l))
		}
	}
	return n
}

func (m *SetSmartContractDeployersProposal_Deployer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainReferenceID)
	if l > 0 {
		n += 1 + l + sovSetSmartContractDeployersProposal(uint64(l))
	}
	l = len(m.ContractAddress)
	if l > 0 {
		n += 1 + l + sovSetSmartContractDeployersProposal(uint64(l))
	}
	return n
}

func sovSetSmartContractDeployersProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSetSmartContractDeployersProposal(x uint64) (n int) {
	return sovSetSmartContractDeployersProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SetSmartContractDeployersProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSetSmartContractDeployersProposal
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
			return fmt.Errorf("proto: SetSmartContractDeployersProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetSmartContractDeployersProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSetSmartContractDeployersProposal
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
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Summary", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSetSmartContractDeployersProposal
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
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Summary = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deployers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSetSmartContractDeployersProposal
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
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Deployers = append(m.Deployers, SetSmartContractDeployersProposal_Deployer{})
			if err := m.Deployers[len(m.Deployers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSetSmartContractDeployersProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSetSmartContractDeployersProposal
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
func (m *SetSmartContractDeployersProposal_Deployer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSetSmartContractDeployersProposal
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
			return fmt.Errorf("proto: Deployer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Deployer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainReferenceID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSetSmartContractDeployersProposal
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
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainReferenceID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSetSmartContractDeployersProposal
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
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSetSmartContractDeployersProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSetSmartContractDeployersProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSetSmartContractDeployersProposal
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
func skipSetSmartContractDeployersProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSetSmartContractDeployersProposal
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
					return 0, ErrIntOverflowSetSmartContractDeployersProposal
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
					return 0, ErrIntOverflowSetSmartContractDeployersProposal
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
				return 0, ErrInvalidLengthSetSmartContractDeployersProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSetSmartContractDeployersProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSetSmartContractDeployersProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSetSmartContractDeployersProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSetSmartContractDeployersProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSetSmartContractDeployersProposal = fmt.Errorf("proto: unexpected end of group")
)
