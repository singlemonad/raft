// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package raft

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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

type MessageType int32

const (
	MessageType_Vote               MessageType = 0
	MessageType_VoteReply          MessageType = 1
	MessageType_AppendEntries      MessageType = 2
	MessageType_AppendEntriesReply MessageType = 3
	MessageType_ElectionTimeout    MessageType = 4
	MessageType_Puase              MessageType = 5
)

var MessageType_name = map[int32]string{
	0: "Vote",
	1: "VoteReply",
	2: "AppendEntries",
	3: "AppendEntriesReply",
	4: "ElectionTimeout",
	5: "Puase",
}

var MessageType_value = map[string]int32{
	"Vote":               0,
	"VoteReply":          1,
	"AppendEntries":      2,
	"AppendEntriesReply": 3,
	"ElectionTimeout":    4,
	"Puase":              5,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

type Entry struct {
	Term                 int64    `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Body                 *any.Any `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *Entry) GetBody() *any.Any {
	if m != nil {
		return m.Body
	}
	return nil
}

type Message struct {
	Type                 MessageType `protobuf:"varint,1,opt,name=type,proto3,enum=raft.MessageType" json:"type,omitempty"`
	To                   int64       `protobuf:"varint,2,opt,name=to,proto3" json:"to,omitempty"`
	From                 int64       `protobuf:"varint,3,opt,name=from,proto3" json:"from,omitempty"`
	Term                 int64       `protobuf:"varint,4,opt,name=term,proto3" json:"term,omitempty"`
	LogIndex             int64       `protobuf:"varint,5,opt,name=logIndex,proto3" json:"logIndex,omitempty"`
	LogTerm              int64       `protobuf:"varint,6,opt,name=logTerm,proto3" json:"logTerm,omitempty"`
	CommitIndex          int64       `protobuf:"varint,7,opt,name=commitIndex,proto3" json:"commitIndex,omitempty"`
	Entries              []*Entry    `protobuf:"bytes,8,rep,name=entries,proto3" json:"entries,omitempty"`
	Reject               bool        `protobuf:"varint,9,opt,name=reject,proto3" json:"reject,omitempty"`
	PasueTime            int64       `protobuf:"varint,10,opt,name=pasueTime,proto3" json:"pasueTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetType() MessageType {
	if m != nil {
		return m.Type
	}
	return MessageType_Vote
}

func (m *Message) GetTo() int64 {
	if m != nil {
		return m.To
	}
	return 0
}

func (m *Message) GetFrom() int64 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *Message) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *Message) GetLogIndex() int64 {
	if m != nil {
		return m.LogIndex
	}
	return 0
}

func (m *Message) GetLogTerm() int64 {
	if m != nil {
		return m.LogTerm
	}
	return 0
}

func (m *Message) GetCommitIndex() int64 {
	if m != nil {
		return m.CommitIndex
	}
	return 0
}

func (m *Message) GetEntries() []*Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *Message) GetReject() bool {
	if m != nil {
		return m.Reject
	}
	return false
}

func (m *Message) GetPasueTime() int64 {
	if m != nil {
		return m.PasueTime
	}
	return 0
}

func init() {
	proto.RegisterEnum("raft.MessageType", MessageType_name, MessageType_value)
	proto.RegisterType((*Entry)(nil), "raft.Entry")
	proto.RegisterType((*Message)(nil), "raft.Message")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 350 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x51, 0x4d, 0x6b, 0xab, 0x50,
	0x10, 0x7d, 0x7e, 0x45, 0x1d, 0x49, 0x9e, 0x99, 0xf7, 0x08, 0xf7, 0x85, 0xb7, 0x90, 0x40, 0x40,
	0xba, 0x30, 0x90, 0xfe, 0x82, 0x2c, 0x5c, 0x74, 0x51, 0x28, 0x12, 0xba, 0x37, 0xc9, 0x44, 0x2c,
	0xea, 0x58, 0xbd, 0x42, 0xfd, 0x31, 0xfd, 0xaf, 0xc5, 0x6b, 0x92, 0xa6, 0xbb, 0x99, 0x73, 0xce,
	0x3d, 0xe7, 0x30, 0x17, 0xa6, 0x25, 0xb5, 0x6d, 0x9a, 0x51, 0x54, 0x37, 0x2c, 0x19, 0xcd, 0x26,
	0x3d, 0xcb, 0xe5, 0xbf, 0x8c, 0x39, 0x2b, 0x68, 0xa3, 0xb0, 0x43, 0x77, 0xde, 0xa4, 0x55, 0x3f,
	0x0a, 0x56, 0x31, 0x58, 0x71, 0x25, 0x9b, 0x1e, 0x11, 0x4c, 0x49, 0x4d, 0x29, 0xf4, 0x40, 0x0b,
	0x8d, 0x44, 0xcd, 0x18, 0x82, 0x79, 0xe0, 0x53, 0x2f, 0x8c, 0x40, 0x0b, 0xbd, 0xed, 0xdf, 0x68,
	0xb4, 0x89, 0xae, 0x36, 0xd1, 0xae, 0xea, 0x13, 0xa5, 0x58, 0x7d, 0xea, 0x60, 0x3f, 0x8f, 0xc9,
	0xb8, 0x06, 0x53, 0xf6, 0x35, 0x09, 0x2d, 0xd0, 0xc2, 0xd9, 0x76, 0x1e, 0x0d, 0x15, 0xa2, 0x0b,
	0xb9, 0xef, 0x6b, 0x4a, 0x14, 0x8d, 0x33, 0xd0, 0x25, 0x5f, 0xe2, 0x74, 0xc9, 0x43, 0x81, 0x73,
	0xc3, 0xa5, 0x0a, 0x33, 0x12, 0x35, 0xdf, 0x4a, 0x99, 0x77, 0xa5, 0x96, 0xe0, 0x14, 0x9c, 0x3d,
	0x55, 0x27, 0xfa, 0x10, 0x96, 0xc2, 0x6f, 0x3b, 0x0a, 0xb0, 0x0b, 0xce, 0xf6, 0xc3, 0x93, 0x89,
	0xa2, 0xae, 0x2b, 0x06, 0xe0, 0x1d, 0xb9, 0x2c, 0x73, 0x39, 0x3e, 0xb4, 0x15, 0x7b, 0x0f, 0xe1,
	0x1a, 0x6c, 0xaa, 0x64, 0x93, 0x53, 0x2b, 0x9c, 0xc0, 0x08, 0xbd, 0xad, 0x37, 0x36, 0x57, 0xe7,
	0x49, 0xae, 0x1c, 0x2e, 0x60, 0xd2, 0xd0, 0x1b, 0x1d, 0xa5, 0x70, 0x03, 0x2d, 0x74, 0x92, 0xcb,
	0x86, 0xff, 0xc1, 0xad, 0xd3, 0xb6, 0xa3, 0x7d, 0x5e, 0x92, 0x00, 0x65, 0xff, 0x0d, 0x3c, 0xbc,
	0x83, 0x77, 0x77, 0x01, 0x74, 0xc0, 0x7c, 0x65, 0x49, 0xfe, 0x2f, 0x9c, 0x82, 0x3b, 0x4c, 0x09,
	0xd5, 0x45, 0xef, 0x6b, 0x38, 0x87, 0xe9, 0xae, 0xae, 0xa9, 0x3a, 0xc5, 0x63, 0x9c, 0xaf, 0xe3,
	0x02, 0xf0, 0x07, 0x34, 0x4a, 0x0d, 0xfc, 0x03, 0xbf, 0xe3, 0x82, 0x8e, 0x32, 0xe7, 0x6a, 0x88,
	0xe0, 0x4e, 0xfa, 0x26, 0xba, 0x60, 0xbd, 0x74, 0x69, 0x4b, 0xbe, 0x75, 0x98, 0xa8, 0x6f, 0x7a,
	0xfc, 0x0a, 0x00, 0x00, 0xff, 0xff, 0xdc, 0xd3, 0x6f, 0x70, 0x12, 0x02, 0x00, 0x00,
}
