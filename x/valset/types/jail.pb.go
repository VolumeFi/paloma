// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: palomachain/paloma/valset/jail.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type JailRecord struct {
	// Address of the validator being jailed
	Address github_com_cosmos_cosmos_sdk_types.ValAddress `protobuf:"bytes,1,opt,name=address,proto3,casttype=github.com/cosmos/cosmos-sdk/types.ValAddress" json:"address,omitempty"`
	// Duration of jailing
	Duration time.Duration `protobuf:"bytes,2,opt,name=duration,proto3,stdduration" json:"duration"`
	// Timestamp of when the validator was jailed
	JailedAt time.Time `protobuf:"bytes,3,opt,name=jailedAt,proto3,stdtime" json:"jailedAt"`
}

func (m *JailRecord) Reset()         { *m = JailRecord{} }
func (m *JailRecord) String() string { return proto.CompactTextString(m) }
func (*JailRecord) ProtoMessage()    {}
func (*JailRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_7d6891913aadee6b, []int{0}
}
func (m *JailRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *JailRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_JailRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *JailRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JailRecord.Merge(m, src)
}
func (m *JailRecord) XXX_Size() int {
	return m.Size()
}
func (m *JailRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_JailRecord.DiscardUnknown(m)
}

var xxx_messageInfo_JailRecord proto.InternalMessageInfo

func (m *JailRecord) GetAddress() github_com_cosmos_cosmos_sdk_types.ValAddress {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *JailRecord) GetDuration() time.Duration {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *JailRecord) GetJailedAt() time.Time {
	if m != nil {
		return m.JailedAt
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*JailRecord)(nil), "palomachain.paloma.valset.JailRecord")
}

func init() {
	proto.RegisterFile("palomachain/paloma/valset/jail.proto", fileDescriptor_7d6891913aadee6b)
}

var fileDescriptor_7d6891913aadee6b = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0x3f, 0x4f, 0x02, 0x31,
	0x14, 0xbf, 0x6a, 0xa2, 0xa4, 0x3a, 0x5d, 0x1c, 0x80, 0xa1, 0x47, 0x8c, 0x03, 0x0b, 0x6d, 0xc4,
	0x0f, 0xa0, 0x10, 0x17, 0x75, 0x23, 0xc6, 0xc1, 0xc5, 0x94, 0x6b, 0x3d, 0xaa, 0x3d, 0xde, 0xe5,
	0x5a, 0x8c, 0x7e, 0x0b, 0x46, 0x3f, 0x12, 0x23, 0xa3, 0x83, 0x41, 0x03, 0xdf, 0xc2, 0xc9, 0xd0,
	0x5e, 0x0d, 0x51, 0xa7, 0xbe, 0x97, 0xdf, 0x9f, 0xf7, 0x7e, 0x7d, 0xf8, 0xa8, 0xe0, 0x1a, 0x72,
	0x9e, 0x8e, 0xb8, 0x1a, 0x33, 0x5f, 0xb3, 0x27, 0xae, 0x8d, 0xb4, 0xec, 0x81, 0x2b, 0x4d, 0x8b,
	0x12, 0x2c, 0xc4, 0x8d, 0x0d, 0x16, 0xf5, 0x35, 0xf5, 0xac, 0xe6, 0x41, 0x06, 0x19, 0x38, 0x16,
	0x5b, 0x57, 0x5e, 0xd0, 0x6c, 0xa4, 0x60, 0x72, 0x30, 0x77, 0x1e, 0xf0, 0x4d, 0x05, 0x91, 0x0c,
	0x20, 0xd3, 0x92, 0xb9, 0x6e, 0x38, 0xb9, 0x67, 0x62, 0x52, 0x72, 0xab, 0x60, 0x5c, 0xe1, 0xc9,
	0x6f, 0xdc, 0xaa, 0x5c, 0x1a, 0xcb, 0xf3, 0xc2, 0x13, 0x0e, 0xdf, 0x11, 0xc6, 0x97, 0x5c, 0xe9,
	0x81, 0x4c, 0xa1, 0x14, 0xf1, 0x15, 0xde, 0xe5, 0x42, 0x94, 0xd2, 0x98, 0x3a, 0x6a, 0xa1, 0xf6,
	0x7e, 0xff, 0xf8, 0x6b, 0x91, 0x74, 0x32, 0x65, 0x47, 0x93, 0x21, 0x4d, 0x21, 0xaf, 0xa6, 0x57,
	0x4f, 0xc7, 0x88, 0x47, 0x66, 0x5f, 0x0a, 0x69, 0xe8, 0x0d, 0xd7, 0x3d, 0x2f, 0x1c, 0x04, 0x87,
	0xf8, 0x14, 0xd7, 0xc2, 0x3a, 0xf5, 0xad, 0x16, 0x6a, 0xef, 0x75, 0x1b, 0xd4, 0xef, 0x43, 0xc3,
	0x3e, 0xf4, 0xbc, 0x22, 0xf4, 0x6b, 0xb3, 0x45, 0x12, 0xbd, 0x7e, 0x24, 0x68, 0xf0, 0x23, 0x8a,
	0xcf, 0x70, 0x6d, 0xfd, 0x6f, 0x52, 0xf4, 0x6c, 0x7d, 0xdb, 0x19, 0x34, 0xff, 0x18, 0x5c, 0x87,
	0x40, 0xde, 0x61, 0xea, 0x1c, 0x82, 0xaa, 0x7f, 0x31, 0x5b, 0x12, 0x34, 0x5f, 0x12, 0xf4, 0xb9,
	0x24, 0x68, 0xba, 0x22, 0xd1, 0x7c, 0x45, 0xa2, 0xb7, 0x15, 0x89, 0x6e, 0xd9, 0x46, 0xa8, 0xff,
	0xce, 0xd6, 0x65, 0xcf, 0xe1, 0x76, 0x2e, 0xe1, 0x70, 0xc7, 0x8d, 0x3c, 0xf9, 0x0e, 0x00, 0x00,
	0xff, 0xff, 0x10, 0xaa, 0xc8, 0x80, 0xe5, 0x01, 0x00, 0x00,
}

func (m *JailRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *JailRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *JailRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.JailedAt, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.JailedAt):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintJail(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x1a
	n2, err2 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.Duration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.Duration):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintJail(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintJail(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintJail(dAtA []byte, offset int, v uint64) int {
	offset -= sovJail(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *JailRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovJail(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.Duration)
	n += 1 + l + sovJail(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.JailedAt)
	n += 1 + l + sovJail(uint64(l))
	return n
}

func sovJail(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozJail(x uint64) (n int) {
	return sovJail(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *JailRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowJail
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
			return fmt.Errorf("proto: JailRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: JailRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJail
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
				return ErrInvalidLengthJail
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthJail
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJail
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
				return ErrInvalidLengthJail
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthJail
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.Duration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JailedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowJail
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
				return ErrInvalidLengthJail
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthJail
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.JailedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipJail(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthJail
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
func skipJail(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowJail
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
					return 0, ErrIntOverflowJail
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
					return 0, ErrIntOverflowJail
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
				return 0, ErrInvalidLengthJail
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupJail
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthJail
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthJail        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowJail          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupJail = fmt.Errorf("proto: unexpected end of group")
)
