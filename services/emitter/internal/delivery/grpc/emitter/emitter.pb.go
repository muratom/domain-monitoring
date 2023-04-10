// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.6.1
// source: emitter/emitter.proto

package emitter

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetDNSRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fqdn string `protobuf:"bytes,1,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
}

func (x *GetDNSRequest) Reset() {
	*x = GetDNSRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_emitter_emitter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDNSRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDNSRequest) ProtoMessage() {}

func (x *GetDNSRequest) ProtoReflect() protoreflect.Message {
	mi := &file_emitter_emitter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDNSRequest.ProtoReflect.Descriptor instead.
func (*GetDNSRequest) Descriptor() ([]byte, []int) {
	return file_emitter_emitter_proto_rawDescGZIP(), []int{0}
}

func (x *GetDNSRequest) GetFqdn() string {
	if x != nil {
		return x.Fqdn
	}
	return ""
}

type MX struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Pref uint64 `protobuf:"varint,2,opt,name=pref,proto3" json:"pref,omitempty"`
}

func (x *MX) Reset() {
	*x = MX{}
	if protoimpl.UnsafeEnabled {
		mi := &file_emitter_emitter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MX) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MX) ProtoMessage() {}

func (x *MX) ProtoReflect() protoreflect.Message {
	mi := &file_emitter_emitter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MX.ProtoReflect.Descriptor instead.
func (*MX) Descriptor() ([]byte, []int) {
	return file_emitter_emitter_proto_rawDescGZIP(), []int{1}
}

func (x *MX) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *MX) GetPref() uint64 {
	if x != nil {
		return x.Pref
	}
	return 0
}

type NS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *NS) Reset() {
	*x = NS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_emitter_emitter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NS) ProtoMessage() {}

func (x *NS) ProtoReflect() protoreflect.Message {
	mi := &file_emitter_emitter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NS.ProtoReflect.Descriptor instead.
func (*NS) Descriptor() ([]byte, []int) {
	return file_emitter_emitter_proto_rawDescGZIP(), []int{2}
}

func (x *NS) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

type SRV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Target   string `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`
	Port     uint64 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Priority uint64 `protobuf:"varint,3,opt,name=priority,proto3" json:"priority,omitempty"`
	Weight   uint64 `protobuf:"varint,4,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *SRV) Reset() {
	*x = SRV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_emitter_emitter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SRV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SRV) ProtoMessage() {}

func (x *SRV) ProtoReflect() protoreflect.Message {
	mi := &file_emitter_emitter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SRV.ProtoReflect.Descriptor instead.
func (*SRV) Descriptor() ([]byte, []int) {
	return file_emitter_emitter_proto_rawDescGZIP(), []int{3}
}

func (x *SRV) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *SRV) GetPort() uint64 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *SRV) GetPriority() uint64 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *SRV) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

type ResourceRecords struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A     []string `protobuf:"bytes,1,rep,name=a,proto3" json:"a,omitempty"`
	Aaaa  []string `protobuf:"bytes,2,rep,name=aaaa,proto3" json:"aaaa,omitempty"`
	Cname string   `protobuf:"bytes,3,opt,name=cname,proto3" json:"cname,omitempty"`
	Mx    []*MX    `protobuf:"bytes,4,rep,name=mx,proto3" json:"mx,omitempty"`
	Ns    []*NS    `protobuf:"bytes,5,rep,name=ns,proto3" json:"ns,omitempty"`
	Srv   []*SRV   `protobuf:"bytes,6,rep,name=srv,proto3" json:"srv,omitempty"`
	Txt   []string `protobuf:"bytes,7,rep,name=txt,proto3" json:"txt,omitempty"`
}

func (x *ResourceRecords) Reset() {
	*x = ResourceRecords{}
	if protoimpl.UnsafeEnabled {
		mi := &file_emitter_emitter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceRecords) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceRecords) ProtoMessage() {}

func (x *ResourceRecords) ProtoReflect() protoreflect.Message {
	mi := &file_emitter_emitter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceRecords.ProtoReflect.Descriptor instead.
func (*ResourceRecords) Descriptor() ([]byte, []int) {
	return file_emitter_emitter_proto_rawDescGZIP(), []int{4}
}

func (x *ResourceRecords) GetA() []string {
	if x != nil {
		return x.A
	}
	return nil
}

func (x *ResourceRecords) GetAaaa() []string {
	if x != nil {
		return x.Aaaa
	}
	return nil
}

func (x *ResourceRecords) GetCname() string {
	if x != nil {
		return x.Cname
	}
	return ""
}

func (x *ResourceRecords) GetMx() []*MX {
	if x != nil {
		return x.Mx
	}
	return nil
}

func (x *ResourceRecords) GetNs() []*NS {
	if x != nil {
		return x.Ns
	}
	return nil
}

func (x *ResourceRecords) GetSrv() []*SRV {
	if x != nil {
		return x.Srv
	}
	return nil
}

func (x *ResourceRecords) GetTxt() []string {
	if x != nil {
		return x.Txt
	}
	return nil
}

type GetWhoisRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fqdn string `protobuf:"bytes,1,opt,name=fqdn,proto3" json:"fqdn,omitempty"`
}

func (x *GetWhoisRequest) Reset() {
	*x = GetWhoisRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_emitter_emitter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWhoisRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWhoisRequest) ProtoMessage() {}

func (x *GetWhoisRequest) ProtoReflect() protoreflect.Message {
	mi := &file_emitter_emitter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWhoisRequest.ProtoReflect.Descriptor instead.
func (*GetWhoisRequest) Descriptor() ([]byte, []int) {
	return file_emitter_emitter_proto_rawDescGZIP(), []int{5}
}

func (x *GetWhoisRequest) GetFqdn() string {
	if x != nil {
		return x.Fqdn
	}
	return ""
}

type WhoisRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DomainName  string               `protobuf:"bytes,1,opt,name=domain_name,json=domainName,proto3" json:"domain_name,omitempty"`
	NameServers []string             `protobuf:"bytes,2,rep,name=name_servers,json=nameServers,proto3" json:"name_servers,omitempty"`
	Created     *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created,proto3" json:"created,omitempty"`
	PaidTill    *timestamp.Timestamp `protobuf:"bytes,4,opt,name=paid_till,json=paidTill,proto3" json:"paid_till,omitempty"`
}

func (x *WhoisRecord) Reset() {
	*x = WhoisRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_emitter_emitter_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WhoisRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WhoisRecord) ProtoMessage() {}

func (x *WhoisRecord) ProtoReflect() protoreflect.Message {
	mi := &file_emitter_emitter_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WhoisRecord.ProtoReflect.Descriptor instead.
func (*WhoisRecord) Descriptor() ([]byte, []int) {
	return file_emitter_emitter_proto_rawDescGZIP(), []int{6}
}

func (x *WhoisRecord) GetDomainName() string {
	if x != nil {
		return x.DomainName
	}
	return ""
}

func (x *WhoisRecord) GetNameServers() []string {
	if x != nil {
		return x.NameServers
	}
	return nil
}

func (x *WhoisRecord) GetCreated() *timestamp.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *WhoisRecord) GetPaidTill() *timestamp.Timestamp {
	if x != nil {
		return x.PaidTill
	}
	return nil
}

var File_emitter_emitter_proto protoreflect.FileDescriptor

var file_emitter_emitter_proto_rawDesc = []byte{
	0x0a, 0x15, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x23, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x44, 0x4e, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x71, 0x64, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x66, 0x71, 0x64, 0x6e, 0x22, 0x2c, 0x0a, 0x02, 0x4d, 0x58, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x72, 0x65, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x70, 0x72,
	0x65, 0x66, 0x22, 0x18, 0x0a, 0x02, 0x4e, 0x53, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0x65, 0x0a, 0x03,
	0x53, 0x52, 0x56, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x77,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x77, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x22, 0xb5, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x0c, 0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x01, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x61, 0x61, 0x61, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x61, 0x61, 0x61, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1b, 0x0a, 0x02, 0x6d, 0x78, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x65, 0x6d,
	0x69, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x58, 0x52, 0x02, 0x6d, 0x78, 0x12, 0x1b, 0x0a, 0x02,
	0x6e, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x65, 0x6d, 0x69, 0x74, 0x74,
	0x65, 0x72, 0x2e, 0x4e, 0x53, 0x52, 0x02, 0x6e, 0x73, 0x12, 0x1e, 0x0a, 0x03, 0x73, 0x72, 0x76,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72,
	0x2e, 0x53, 0x52, 0x56, 0x52, 0x03, 0x73, 0x72, 0x76, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x78, 0x74,
	0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x74, 0x78, 0x74, 0x22, 0x25, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x57, 0x68, 0x6f, 0x69, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x66, 0x71, 0x64, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x71,
	0x64, 0x6e, 0x22, 0xc0, 0x01, 0x0a, 0x0b, 0x57, 0x68, 0x6f, 0x69, 0x73, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x61, 0x6d, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x37, 0x0a, 0x09,
	0x70, 0x61, 0x69, 0x64, 0x5f, 0x74, 0x69, 0x6c, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x70, 0x61, 0x69,
	0x64, 0x54, 0x69, 0x6c, 0x6c, 0x32, 0xac, 0x01, 0x0a, 0x0e, 0x45, 0x6d, 0x69, 0x74, 0x74, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x44,
	0x4e, 0x53, 0x12, 0x16, 0x2e, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74,
	0x44, 0x4e, 0x53, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x65, 0x6d, 0x69,
	0x74, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x22, 0x0f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x2f, 0x76,
	0x31, 0x2f, 0x64, 0x6e, 0x73, 0x12, 0x4d, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x57, 0x68, 0x6f, 0x69,
	0x73, 0x12, 0x18, 0x2e, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x57,
	0x68, 0x6f, 0x69, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x65, 0x6d,
	0x69, 0x74, 0x74, 0x65, 0x72, 0x2e, 0x57, 0x68, 0x6f, 0x69, 0x73, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x77,
	0x68, 0x6f, 0x69, 0x73, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x72, 0x61, 0x74, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x2d, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72, 0x2f, 0x64, 0x65,
	0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x6d, 0x69, 0x74,
	0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_emitter_emitter_proto_rawDescOnce sync.Once
	file_emitter_emitter_proto_rawDescData = file_emitter_emitter_proto_rawDesc
)

func file_emitter_emitter_proto_rawDescGZIP() []byte {
	file_emitter_emitter_proto_rawDescOnce.Do(func() {
		file_emitter_emitter_proto_rawDescData = protoimpl.X.CompressGZIP(file_emitter_emitter_proto_rawDescData)
	})
	return file_emitter_emitter_proto_rawDescData
}

var file_emitter_emitter_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_emitter_emitter_proto_goTypes = []interface{}{
	(*GetDNSRequest)(nil),       // 0: emitter.GetDNSRequest
	(*MX)(nil),                  // 1: emitter.MX
	(*NS)(nil),                  // 2: emitter.NS
	(*SRV)(nil),                 // 3: emitter.SRV
	(*ResourceRecords)(nil),     // 4: emitter.ResourceRecords
	(*GetWhoisRequest)(nil),     // 5: emitter.GetWhoisRequest
	(*WhoisRecord)(nil),         // 6: emitter.WhoisRecord
	(*timestamp.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_emitter_emitter_proto_depIdxs = []int32{
	1, // 0: emitter.ResourceRecords.mx:type_name -> emitter.MX
	2, // 1: emitter.ResourceRecords.ns:type_name -> emitter.NS
	3, // 2: emitter.ResourceRecords.srv:type_name -> emitter.SRV
	7, // 3: emitter.WhoisRecord.created:type_name -> google.protobuf.Timestamp
	7, // 4: emitter.WhoisRecord.paid_till:type_name -> google.protobuf.Timestamp
	0, // 5: emitter.EmitterService.GetDNS:input_type -> emitter.GetDNSRequest
	5, // 6: emitter.EmitterService.GetWhois:input_type -> emitter.GetWhoisRequest
	4, // 7: emitter.EmitterService.GetDNS:output_type -> emitter.ResourceRecords
	6, // 8: emitter.EmitterService.GetWhois:output_type -> emitter.WhoisRecord
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_emitter_emitter_proto_init() }
func file_emitter_emitter_proto_init() {
	if File_emitter_emitter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_emitter_emitter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDNSRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_emitter_emitter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MX); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_emitter_emitter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NS); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_emitter_emitter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SRV); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_emitter_emitter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceRecords); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_emitter_emitter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWhoisRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_emitter_emitter_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WhoisRecord); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_emitter_emitter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_emitter_emitter_proto_goTypes,
		DependencyIndexes: file_emitter_emitter_proto_depIdxs,
		MessageInfos:      file_emitter_emitter_proto_msgTypes,
	}.Build()
	File_emitter_emitter_proto = out.File
	file_emitter_emitter_proto_rawDesc = nil
	file_emitter_emitter_proto_goTypes = nil
	file_emitter_emitter_proto_depIdxs = nil
}
