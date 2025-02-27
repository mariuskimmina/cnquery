// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: resources.proto

package resources

import (
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

type Schema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resources map[string]*ResourceInfo `protobuf:"bytes,3,rep,name=resources,proto3" json:"resources,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Schema) Reset() {
	*x = Schema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Schema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schema) ProtoMessage() {}

func (x *Schema) ProtoReflect() protoreflect.Message {
	mi := &file_resources_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Schema.ProtoReflect.Descriptor instead.
func (*Schema) Descriptor() ([]byte, []int) {
	return file_resources_proto_rawDescGZIP(), []int{0}
}

func (x *Schema) GetResources() map[string]*ResourceInfo {
	if x != nil {
		return x.Resources
	}
	return nil
}

type ResourceID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ResourceID) Reset() {
	*x = ResourceID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceID) ProtoMessage() {}

func (x *ResourceID) ProtoReflect() protoreflect.Message {
	mi := &file_resources_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceID.ProtoReflect.Descriptor instead.
func (*ResourceID) Descriptor() ([]byte, []int) {
	return file_resources_proto_rawDescGZIP(), []int{1}
}

func (x *ResourceID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ResourceID) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type TypedArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type     string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Optional bool   `protobuf:"varint,3,opt,name=optional,proto3" json:"optional,omitempty"`
}

func (x *TypedArg) Reset() {
	*x = TypedArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TypedArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypedArg) ProtoMessage() {}

func (x *TypedArg) ProtoReflect() protoreflect.Message {
	mi := &file_resources_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypedArg.ProtoReflect.Descriptor instead.
func (*TypedArg) Descriptor() ([]byte, []int) {
	return file_resources_proto_rawDescGZIP(), []int{2}
}

func (x *TypedArg) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TypedArg) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *TypedArg) GetOptional() bool {
	if x != nil {
		return x.Optional
	}
	return false
}

type Init struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Args []*TypedArg `protobuf:"bytes,1,rep,name=args,proto3" json:"args,omitempty"`
}

func (x *Init) Reset() {
	*x = Init{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Init) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Init) ProtoMessage() {}

func (x *Init) ProtoReflect() protoreflect.Message {
	mi := &file_resources_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Init.ProtoReflect.Descriptor instead.
func (*Init) Descriptor() ([]byte, []int) {
	return file_resources_proto_rawDescGZIP(), []int{3}
}

func (x *Init) GetArgs() []*TypedArg {
	if x != nil {
		return x.Args
	}
	return nil
}

type ResourceInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name             string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Fields           map[string]*Field `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Init             *Init             `protobuf:"bytes,20,opt,name=init,proto3" json:"init,omitempty"`
	ListType         string            `protobuf:"bytes,21,opt,name=list_type,json=listType,proto3" json:"list_type,omitempty"`
	Title            string            `protobuf:"bytes,22,opt,name=title,proto3" json:"title,omitempty"`
	Desc             string            `protobuf:"bytes,23,opt,name=desc,proto3" json:"desc,omitempty"`
	Private          bool              `protobuf:"varint,24,opt,name=private,proto3" json:"private,omitempty"`
	IsExtension      bool              `protobuf:"varint,28,opt,name=is_extension,json=isExtension,proto3" json:"is_extension,omitempty"`
	MinMondooVersion string            `protobuf:"bytes,25,opt,name=min_mondoo_version,json=minMondooVersion,proto3" json:"min_mondoo_version,omitempty"`
	Defaults         string            `protobuf:"bytes,26,opt,name=defaults,proto3" json:"defaults,omitempty"`
	Provider         string            `protobuf:"bytes,27,opt,name=provider,proto3" json:"provider,omitempty"`
}

func (x *ResourceInfo) Reset() {
	*x = ResourceInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceInfo) ProtoMessage() {}

func (x *ResourceInfo) ProtoReflect() protoreflect.Message {
	mi := &file_resources_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceInfo.ProtoReflect.Descriptor instead.
func (*ResourceInfo) Descriptor() ([]byte, []int) {
	return file_resources_proto_rawDescGZIP(), []int{4}
}

func (x *ResourceInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ResourceInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ResourceInfo) GetFields() map[string]*Field {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ResourceInfo) GetInit() *Init {
	if x != nil {
		return x.Init
	}
	return nil
}

func (x *ResourceInfo) GetListType() string {
	if x != nil {
		return x.ListType
	}
	return ""
}

func (x *ResourceInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ResourceInfo) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *ResourceInfo) GetPrivate() bool {
	if x != nil {
		return x.Private
	}
	return false
}

func (x *ResourceInfo) GetIsExtension() bool {
	if x != nil {
		return x.IsExtension
	}
	return false
}

func (x *ResourceInfo) GetMinMondooVersion() string {
	if x != nil {
		return x.MinMondooVersion
	}
	return ""
}

func (x *ResourceInfo) GetDefaults() string {
	if x != nil {
		return x.Defaults
	}
	return ""
}

func (x *ResourceInfo) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name               string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type               string   `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	IsMandatory        bool     `protobuf:"varint,3,opt,name=is_mandatory,json=isMandatory,proto3" json:"is_mandatory,omitempty"`
	Refs               []string `protobuf:"bytes,4,rep,name=refs,proto3" json:"refs,omitempty"`
	Title              string   `protobuf:"bytes,20,opt,name=title,proto3" json:"title,omitempty"`
	Desc               string   `protobuf:"bytes,21,opt,name=desc,proto3" json:"desc,omitempty"`
	IsPrivate          bool     `protobuf:"varint,22,opt,name=is_private,json=isPrivate,proto3" json:"is_private,omitempty"`
	MinMondooVersion   string   `protobuf:"bytes,23,opt,name=min_mondoo_version,json=minMondooVersion,proto3" json:"min_mondoo_version,omitempty"`
	IsImplicitResource bool     `protobuf:"varint,24,opt,name=is_implicit_resource,json=isImplicitResource,proto3" json:"is_implicit_resource,omitempty"`
	IsEmbedded         bool     `protobuf:"varint,25,opt,name=is_embedded,json=isEmbedded,proto3" json:"is_embedded,omitempty"`
	Provider           string   `protobuf:"bytes,27,opt,name=provider,proto3" json:"provider,omitempty"`
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_resources_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_resources_proto_rawDescGZIP(), []int{5}
}

func (x *Field) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Field) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Field) GetIsMandatory() bool {
	if x != nil {
		return x.IsMandatory
	}
	return false
}

func (x *Field) GetRefs() []string {
	if x != nil {
		return x.Refs
	}
	return nil
}

func (x *Field) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Field) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *Field) GetIsPrivate() bool {
	if x != nil {
		return x.IsPrivate
	}
	return false
}

func (x *Field) GetMinMondooVersion() string {
	if x != nil {
		return x.MinMondooVersion
	}
	return ""
}

func (x *Field) GetIsImplicitResource() bool {
	if x != nil {
		return x.IsImplicitResource
	}
	return false
}

func (x *Field) GetIsEmbedded() bool {
	if x != nil {
		return x.IsEmbedded
	}
	return false
}

func (x *Field) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

var File_resources_proto protoreflect.FileDescriptor

var file_resources_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x6d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x22, 0xad, 0x01, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x45,
	0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x27, 0x2e, 0x6d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x1a, 0x5c, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x34, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6d, 0x6f, 0x6e, 0x64, 0x6f,
	0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x30, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49,
	0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4e, 0x0a, 0x08, 0x54, 0x79, 0x70, 0x65, 0x64, 0x41, 0x72,
	0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x22, 0x36, 0x0a, 0x04, 0x49, 0x6e, 0x69, 0x74, 0x12, 0x2e, 0x0a,
	0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6d, 0x6f,
	0x6e, 0x64, 0x6f, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x54,
	0x79, 0x70, 0x65, 0x64, 0x41, 0x72, 0x67, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x22, 0xe0, 0x03,
	0x0a, 0x0c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x42, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x2a, 0x0a, 0x04, 0x69, 0x6e, 0x69, 0x74, 0x18, 0x14,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x52, 0x04, 0x69, 0x6e,
	0x69, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x15, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x17, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x18, 0x18, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70, 0x72, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x45, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x12, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f,
	0x6e, 0x64, 0x6f, 0x6f, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x19, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x10, 0x6d, 0x69, 0x6e, 0x4d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73,
	0x18, 0x1a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x1b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x1a, 0x52, 0x0a, 0x0b,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d,
	0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0xcc, 0x02, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x4d, 0x61, 0x6e, 0x64,
	0x61, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x66, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x72, 0x65, 0x66, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64,
	0x65, 0x73, 0x63, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x18, 0x16, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x50, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x6d, 0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f,
	0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10,
	0x6d, 0x69, 0x6e, 0x4d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x30, 0x0a, 0x14, 0x69, 0x73, 0x5f, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x5f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x18, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12,
	0x69, 0x73, 0x49, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65,
	0x64, 0x18, 0x19, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x45, 0x6d, 0x62, 0x65, 0x64,
	0x64, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18,
	0x1b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x42,
	0x35, 0x5a, 0x33, 0x67, 0x6f, 0x2e, 0x6d, 0x6f, 0x6e, 0x64, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x6e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2f, 0x76, 0x39, 0x2f, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x73, 0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_proto_rawDescOnce sync.Once
	file_resources_proto_rawDescData = file_resources_proto_rawDesc
)

func file_resources_proto_rawDescGZIP() []byte {
	file_resources_proto_rawDescOnce.Do(func() {
		file_resources_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_proto_rawDescData)
	})
	return file_resources_proto_rawDescData
}

var file_resources_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_resources_proto_goTypes = []interface{}{
	(*Schema)(nil),       // 0: mondoo.resources.Schema
	(*ResourceID)(nil),   // 1: mondoo.resources.ResourceID
	(*TypedArg)(nil),     // 2: mondoo.resources.TypedArg
	(*Init)(nil),         // 3: mondoo.resources.Init
	(*ResourceInfo)(nil), // 4: mondoo.resources.ResourceInfo
	(*Field)(nil),        // 5: mondoo.resources.Field
	nil,                  // 6: mondoo.resources.Schema.ResourcesEntry
	nil,                  // 7: mondoo.resources.ResourceInfo.FieldsEntry
}
var file_resources_proto_depIdxs = []int32{
	6, // 0: mondoo.resources.Schema.resources:type_name -> mondoo.resources.Schema.ResourcesEntry
	2, // 1: mondoo.resources.Init.args:type_name -> mondoo.resources.TypedArg
	7, // 2: mondoo.resources.ResourceInfo.fields:type_name -> mondoo.resources.ResourceInfo.FieldsEntry
	3, // 3: mondoo.resources.ResourceInfo.init:type_name -> mondoo.resources.Init
	4, // 4: mondoo.resources.Schema.ResourcesEntry.value:type_name -> mondoo.resources.ResourceInfo
	5, // 5: mondoo.resources.ResourceInfo.FieldsEntry.value:type_name -> mondoo.resources.Field
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_resources_proto_init() }
func file_resources_proto_init() {
	if File_resources_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Schema); i {
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
		file_resources_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceID); i {
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
		file_resources_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TypedArg); i {
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
		file_resources_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Init); i {
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
		file_resources_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceInfo); i {
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
		file_resources_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
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
			RawDescriptor: file_resources_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_proto_goTypes,
		DependencyIndexes: file_resources_proto_depIdxs,
		MessageInfos:      file_resources_proto_msgTypes,
	}.Build()
	File_resources_proto = out.File
	file_resources_proto_rawDesc = nil
	file_resources_proto_goTypes = nil
	file_resources_proto_depIdxs = nil
}
