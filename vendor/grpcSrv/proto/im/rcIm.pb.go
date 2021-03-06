// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rcIm.proto

package im

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

//发送信息请求数据结构
type Request struct {
	FromUserId           string   `protobuf:"bytes,1,opt,name=fromUserId,proto3" json:"fromUserId,omitempty"`
	ToUids               []string `protobuf:"bytes,2,rep,name=toUids,proto3" json:"toUids,omitempty"`
	Content              string   `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	PushContent          string   `protobuf:"bytes,4,opt,name=pushContent,proto3" json:"pushContent,omitempty"`
	PushData             string   `protobuf:"bytes,5,opt,name=pushData,proto3" json:"pushData,omitempty"`
	MsgType              string   `protobuf:"bytes,6,opt,name=msgType,proto3" json:"msgType,omitempty"`
	ExtraData            string   `protobuf:"bytes,7,opt,name=extraData,proto3" json:"extraData,omitempty"`
	Tags                 []string `protobuf:"bytes,8,rep,name=tags,proto3" json:"tags,omitempty"`
	Duration             int64    `protobuf:"varint,9,opt,name=duration,proto3" json:"duration,omitempty"`
	UserInfo             string   `protobuf:"bytes,10,opt,name=userInfo,proto3" json:"userInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_14c0dd3d013a8d19, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetFromUserId() string {
	if m != nil {
		return m.FromUserId
	}
	return ""
}

func (m *Request) GetToUids() []string {
	if m != nil {
		return m.ToUids
	}
	return nil
}

func (m *Request) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Request) GetPushContent() string {
	if m != nil {
		return m.PushContent
	}
	return ""
}

func (m *Request) GetPushData() string {
	if m != nil {
		return m.PushData
	}
	return ""
}

func (m *Request) GetMsgType() string {
	if m != nil {
		return m.MsgType
	}
	return ""
}

func (m *Request) GetExtraData() string {
	if m != nil {
		return m.ExtraData
	}
	return ""
}

func (m *Request) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Request) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *Request) GetUserInfo() string {
	if m != nil {
		return m.UserInfo
	}
	return ""
}

//发送响应数据结构
type Response struct {
	Code                 int32    `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_14c0dd3d013a8d19, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Response) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "im.Request")
	proto.RegisterType((*Response)(nil), "im.Response")
}

func init() { proto.RegisterFile("rcIm.proto", fileDescriptor_14c0dd3d013a8d19) }

var fileDescriptor_14c0dd3d013a8d19 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0x49, 0xd2, 0xe6, 0xcf, 0x54, 0x50, 0xf6, 0x20, 0x4b, 0x11, 0x09, 0xf5, 0x0f, 0x39,
	0x48, 0x11, 0x7d, 0x03, 0xeb, 0xa5, 0x07, 0x41, 0xa2, 0x7d, 0x80, 0x6d, 0xb2, 0x4d, 0x03, 0x26,
	0x1b, 0x77, 0x26, 0x60, 0x1f, 0xc6, 0x97, 0xf4, 0x09, 0x64, 0x37, 0x69, 0xed, 0x29, 0xde, 0xe6,
	0x37, 0xdf, 0x37, 0x5f, 0x66, 0xc2, 0x02, 0xe8, 0x6c, 0x59, 0xcd, 0x1b, 0xad, 0x48, 0x31, 0xb7,
	0xac, 0x66, 0xdf, 0x2e, 0x04, 0xa9, 0xfc, 0x6c, 0x25, 0x12, 0xbb, 0x04, 0xd8, 0x68, 0x55, 0xad,
	0x50, 0xea, 0x65, 0xce, 0x9d, 0xd8, 0x49, 0xa2, 0xf4, 0xa8, 0xc3, 0xce, 0xc1, 0x27, 0xb5, 0x2a,
	0x73, 0xe4, 0x6e, 0xec, 0x25, 0x51, 0xda, 0x13, 0xe3, 0x10, 0x64, 0xaa, 0x26, 0x59, 0x13, 0xf7,
	0xec, 0xd0, 0x1e, 0x59, 0x0c, 0x93, 0xa6, 0xc5, 0xed, 0xa2, 0x57, 0x47, 0x56, 0x3d, 0x6e, 0xb1,
	0x29, 0x84, 0x06, 0x9f, 0x05, 0x09, 0x3e, 0xb6, 0xf2, 0x81, 0x4d, 0x6e, 0x85, 0xc5, 0xfb, 0xae,
	0x91, 0xdc, 0xef, 0x72, 0x7b, 0x64, 0x17, 0x10, 0xc9, 0x2f, 0xd2, 0xc2, 0x8e, 0x05, 0x56, 0xfb,
	0x6b, 0x30, 0x06, 0x23, 0x12, 0x05, 0xf2, 0xd0, 0x6e, 0x69, 0x6b, 0xf3, 0x9d, 0xbc, 0xd5, 0x82,
	0x4a, 0x55, 0xf3, 0x28, 0x76, 0x12, 0x2f, 0x3d, 0xb0, 0xd1, 0x5a, 0x73, 0x61, 0xbd, 0x51, 0x1c,
	0xba, 0x1d, 0xf6, 0x3c, 0xbb, 0x87, 0x30, 0x95, 0xd8, 0xa8, 0x1a, 0xa5, 0xc9, 0xcd, 0x54, 0x2e,
	0xb9, 0x1b, 0x3b, 0xc9, 0x38, 0xb5, 0x35, 0x3b, 0x03, 0xaf, 0xc2, 0xa2, 0xbf, 0xdb, 0x94, 0x0f,
	0x3f, 0x0e, 0x78, 0x2f, 0x58, 0xb0, 0x6b, 0x08, 0x5e, 0xdb, 0xf5, 0x47, 0x89, 0x5b, 0x36, 0x99,
	0x97, 0xd5, 0xbc, 0xff, 0xcb, 0xd3, 0x93, 0x0e, 0xfa, 0xcc, 0x2b, 0xf0, 0xdf, 0x76, 0x48, 0xb2,
	0x1a, 0x32, 0xdd, 0x42, 0xf4, 0xa4, 0x95, 0xc8, 0x33, 0x81, 0xf4, 0x4f, 0xd8, 0xa2, 0x45, 0x52,
	0x83, 0x61, 0x37, 0x10, 0x2e, 0xb6, 0x82, 0xb4, 0x1a, 0xb6, 0xdd, 0xc1, 0xe9, 0xde, 0xd6, 0x65,
	0xca, 0x01, 0xf7, 0xda, 0xb7, 0x2f, 0xea, 0xf1, 0x37, 0x00, 0x00, 0xff, 0xff, 0x76, 0xe6, 0x95,
	0x5c, 0x5f, 0x02, 0x00, 0x00,
}
