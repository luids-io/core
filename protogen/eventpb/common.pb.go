// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package eventpb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Event_Type int32

const (
	Event_SECURITY Event_Type = 0
)

var Event_Type_name = map[int32]string{
	0: "SECURITY",
}

var Event_Type_value = map[string]int32{
	"SECURITY": 0,
}

func (x Event_Type) String() string {
	return proto.EnumName(Event_Type_name, int32(x))
}

func (Event_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0, 0}
}

type Event_Level int32

const (
	Event_INFO     Event_Level = 0
	Event_LOW      Event_Level = 1
	Event_MEDIUM   Event_Level = 2
	Event_HIGH     Event_Level = 3
	Event_CRITICAL Event_Level = 4
)

var Event_Level_name = map[int32]string{
	0: "INFO",
	1: "LOW",
	2: "MEDIUM",
	3: "HIGH",
	4: "CRITICAL",
}

var Event_Level_value = map[string]int32{
	"INFO":     0,
	"LOW":      1,
	"MEDIUM":   2,
	"HIGH":     3,
	"CRITICAL": 4,
}

func (x Event_Level) String() string {
	return proto.EnumName(Event_Level_name, int32(x))
}

func (Event_Level) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0, 1}
}

type Event_Data_DataEnc int32

const (
	Event_Data_NODATA Event_Data_DataEnc = 0
	Event_Data_JSON   Event_Data_DataEnc = 1
)

var Event_Data_DataEnc_name = map[int32]string{
	0: "NODATA",
	1: "JSON",
}

var Event_Data_DataEnc_value = map[string]int32{
	"NODATA": 0,
	"JSON":   1,
}

func (x Event_Data_DataEnc) String() string {
	return proto.EnumName(Event_Data_DataEnc_name, int32(x))
}

func (Event_Data_DataEnc) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0, 1, 0}
}

type Event struct {
	Id                   string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 Event_Type           `protobuf:"varint,2,opt,name=type,proto3,enum=luids.event.v1.Event_Type" json:"type,omitempty"`
	Code                 int32                `protobuf:"varint,3,opt,name=code,proto3" json:"code,omitempty"`
	Level                Event_Level          `protobuf:"varint,4,opt,name=level,proto3,enum=luids.event.v1.Event_Level" json:"level,omitempty"`
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Source               *Event_Source        `protobuf:"bytes,6,opt,name=source,proto3" json:"source,omitempty"`
	Data                 *Event_Data          `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetType() Event_Type {
	if m != nil {
		return m.Type
	}
	return Event_SECURITY
}

func (m *Event) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Event) GetLevel() Event_Level {
	if m != nil {
		return m.Level
	}
	return Event_INFO
}

func (m *Event) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Event) GetSource() *Event_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *Event) GetData() *Event_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

type Event_Source struct {
	Hostname             string   `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Program              string   `protobuf:"bytes,2,opt,name=program,proto3" json:"program,omitempty"`
	Instance             string   `protobuf:"bytes,3,opt,name=instance,proto3" json:"instance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event_Source) Reset()         { *m = Event_Source{} }
func (m *Event_Source) String() string { return proto.CompactTextString(m) }
func (*Event_Source) ProtoMessage()    {}
func (*Event_Source) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0, 0}
}

func (m *Event_Source) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event_Source.Unmarshal(m, b)
}
func (m *Event_Source) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event_Source.Marshal(b, m, deterministic)
}
func (m *Event_Source) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_Source.Merge(m, src)
}
func (m *Event_Source) XXX_Size() int {
	return xxx_messageInfo_Event_Source.Size(m)
}
func (m *Event_Source) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_Source.DiscardUnknown(m)
}

var xxx_messageInfo_Event_Source proto.InternalMessageInfo

func (m *Event_Source) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *Event_Source) GetProgram() string {
	if m != nil {
		return m.Program
	}
	return ""
}

func (m *Event_Source) GetInstance() string {
	if m != nil {
		return m.Instance
	}
	return ""
}

type Event_Data struct {
	DataEnc              Event_Data_DataEnc `protobuf:"varint,1,opt,name=dataEnc,proto3,enum=luids.event.v1.Event_Data_DataEnc" json:"dataEnc,omitempty"`
	Data                 []byte             `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Event_Data) Reset()         { *m = Event_Data{} }
func (m *Event_Data) String() string { return proto.CompactTextString(m) }
func (*Event_Data) ProtoMessage()    {}
func (*Event_Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0, 1}
}

func (m *Event_Data) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event_Data.Unmarshal(m, b)
}
func (m *Event_Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event_Data.Marshal(b, m, deterministic)
}
func (m *Event_Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_Data.Merge(m, src)
}
func (m *Event_Data) XXX_Size() int {
	return xxx_messageInfo_Event_Data.Size(m)
}
func (m *Event_Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Event_Data proto.InternalMessageInfo

func (m *Event_Data) GetDataEnc() Event_Data_DataEnc {
	if m != nil {
		return m.DataEnc
	}
	return Event_Data_NODATA
}

func (m *Event_Data) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("luids.event.v1.Event_Type", Event_Type_name, Event_Type_value)
	proto.RegisterEnum("luids.event.v1.Event_Level", Event_Level_name, Event_Level_value)
	proto.RegisterEnum("luids.event.v1.Event_Data_DataEnc", Event_Data_DataEnc_name, Event_Data_DataEnc_value)
	proto.RegisterType((*Event)(nil), "luids.event.v1.Event")
	proto.RegisterType((*Event_Source)(nil), "luids.event.v1.Event.Source")
	proto.RegisterType((*Event_Data)(nil), "luids.event.v1.Event.Data")
}

func init() { proto.RegisterFile("common.proto", fileDescriptor_555bd8c177793206) }

var fileDescriptor_555bd8c177793206 = []byte{
	// 457 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x5f, 0xab, 0xda, 0x30,
	0x18, 0xc6, 0xad, 0xd6, 0x56, 0xdf, 0x89, 0x94, 0xb0, 0x8b, 0xd2, 0x0d, 0x26, 0xbd, 0x72, 0x8c,
	0xa5, 0x1c, 0xb7, 0x8b, 0xc1, 0xc6, 0xc0, 0x69, 0xb7, 0xd3, 0xe1, 0x51, 0x88, 0x3d, 0x8c, 0x9d,
	0xbb, 0xd8, 0x66, 0x3d, 0x85, 0x36, 0x29, 0x6d, 0x14, 0xfc, 0xac, 0xfb, 0x32, 0x23, 0xa9, 0x7a,
	0x36, 0x70, 0x37, 0x92, 0xd8, 0xdf, 0xf3, 0x3c, 0xef, 0x9f, 0xc0, 0x28, 0x11, 0x65, 0x29, 0x38,
	0xae, 0x6a, 0x21, 0x05, 0x1a, 0x17, 0xfb, 0x3c, 0x6d, 0x30, 0x3b, 0x30, 0x2e, 0xf1, 0xe1, 0xc6,
	0xfb, 0x98, 0xe5, 0xf2, 0x71, 0xbf, 0xc3, 0x89, 0x28, 0x83, 0x4c, 0x14, 0x94, 0x67, 0x81, 0x06,
	0x77, 0xfb, 0x5f, 0x41, 0x25, 0x8f, 0x15, 0x6b, 0x02, 0x99, 0x97, 0xac, 0x91, 0xb4, 0xac, 0x9e,
	0x4e, 0xad, 0x99, 0xff, 0xdb, 0x84, 0x7e, 0xa8, 0x9c, 0xd0, 0x18, 0xba, 0x79, 0xea, 0x1a, 0x13,
	0x63, 0x3a, 0x24, 0xdd, 0x3c, 0x45, 0x18, 0x4c, 0x25, 0x76, 0xbb, 0x13, 0x63, 0x3a, 0x9e, 0x79,
	0xf8, 0xdf, 0x54, 0xac, 0x45, 0x38, 0x3e, 0x56, 0x8c, 0x68, 0x0e, 0x21, 0x30, 0x13, 0x91, 0x32,
	0xb7, 0x37, 0x31, 0xa6, 0x7d, 0xa2, 0xcf, 0xe8, 0x06, 0xfa, 0x05, 0x3b, 0xb0, 0xc2, 0x35, 0xb5,
	0xc9, 0x8b, 0xeb, 0x26, 0x2b, 0x85, 0x90, 0x96, 0x44, 0x1f, 0x60, 0x78, 0xa9, 0xd1, 0xed, 0x4f,
	0x8c, 0xe9, 0xb3, 0x99, 0x87, 0x33, 0x21, 0xb2, 0x82, 0xe1, 0x73, 0x5b, 0x38, 0x3e, 0x13, 0xe4,
	0x09, 0x46, 0xef, 0xc1, 0x6a, 0xc4, 0xbe, 0x4e, 0x98, 0x6b, 0x69, 0xd9, 0xcb, 0xeb, 0x69, 0x5b,
	0xcd, 0x90, 0x13, 0xab, 0xda, 0x4c, 0xa9, 0xa4, 0xae, 0x7d, 0x8a, 0xba, 0xaa, 0x59, 0x52, 0x49,
	0x89, 0xe6, 0xbc, 0x07, 0xb0, 0x5a, 0x07, 0xe4, 0xc1, 0xe0, 0x51, 0x34, 0x92, 0xd3, 0x92, 0x9d,
	0xc6, 0x76, 0xb9, 0x23, 0x17, 0xec, 0xaa, 0x16, 0x59, 0x4d, 0x4b, 0x3d, 0xbf, 0x21, 0x39, 0x5f,
	0x95, 0x2a, 0xe7, 0x8d, 0xa4, 0x3c, 0x69, 0x47, 0x35, 0x24, 0x97, 0xbb, 0x77, 0x04, 0x53, 0x25,
	0xa1, 0x4f, 0x60, 0xab, 0xac, 0x90, 0x27, 0xda, 0x78, 0x3c, 0xf3, 0xff, 0x5f, 0x96, 0xfe, 0x09,
	0x79, 0x42, 0xce, 0x12, 0xb5, 0x08, 0xdd, 0x91, 0x0a, 0x1e, 0xb5, 0x55, 0xfb, 0xaf, 0xc0, 0x3e,
	0x71, 0x08, 0xc0, 0x5a, 0x6f, 0x96, 0xf3, 0x78, 0xee, 0x74, 0xd0, 0x00, 0xcc, 0xef, 0xdb, 0xcd,
	0xda, 0x31, 0xfc, 0xe7, 0x60, 0xaa, 0x5d, 0xa2, 0x11, 0x0c, 0xb6, 0xe1, 0xe2, 0x9e, 0x44, 0xf1,
	0x4f, 0xa7, 0xe3, 0x7f, 0x86, 0xbe, 0x5e, 0x8e, 0x02, 0xa3, 0xf5, 0xd7, 0x8d, 0xd3, 0x41, 0x36,
	0xf4, 0x56, 0x9b, 0x1f, 0x8e, 0xa1, 0x7c, 0xee, 0xc2, 0x65, 0x74, 0x7f, 0xe7, 0x74, 0xd5, 0xe7,
	0xdb, 0xe8, 0xdb, 0xad, 0xd3, 0x53, 0xfa, 0x05, 0x89, 0xe2, 0x68, 0x31, 0x5f, 0x39, 0xe6, 0x97,
	0x37, 0x0f, 0xaf, 0xff, 0x7a, 0x9c, 0xba, 0x87, 0xb7, 0xb9, 0x08, 0x12, 0x51, 0xb3, 0xf6, 0x8d,
	0x66, 0x8c, 0x07, 0xba, 0xa9, 0x6a, 0xb7, 0xb3, 0xf4, 0x3f, 0xef, 0xfe, 0x04, 0x00, 0x00, 0xff,
	0xff, 0xc7, 0x04, 0x4b, 0xed, 0xee, 0x02, 0x00, 0x00,
}
