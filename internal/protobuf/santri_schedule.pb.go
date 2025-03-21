// santri.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.2
// source: santri_schedule.proto

package proto

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

type SantriSchedule struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	StartPresence string                 `protobuf:"bytes,4,opt,name=startPresence,proto3" json:"startPresence,omitempty"`
	StartTime     string                 `protobuf:"bytes,5,opt,name=startTime,proto3" json:"startTime,omitempty"`
	FinishTime    string                 `protobuf:"bytes,6,opt,name=finishTime,proto3" json:"finishTime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SantriSchedule) Reset() {
	*x = SantriSchedule{}
	mi := &file_santri_schedule_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SantriSchedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SantriSchedule) ProtoMessage() {}

func (x *SantriSchedule) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SantriSchedule.ProtoReflect.Descriptor instead.
func (*SantriSchedule) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{0}
}

func (x *SantriSchedule) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SantriSchedule) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SantriSchedule) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SantriSchedule) GetStartPresence() string {
	if x != nil {
		return x.StartPresence
	}
	return ""
}

func (x *SantriSchedule) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *SantriSchedule) GetFinishTime() string {
	if x != nil {
		return x.FinishTime
	}
	return ""
}

type CreateSantriScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	StartPresence string                 `protobuf:"bytes,3,opt,name=startPresence,proto3" json:"startPresence,omitempty"`
	StartTime     string                 `protobuf:"bytes,4,opt,name=startTime,proto3" json:"startTime,omitempty"`
	FinishTime    string                 `protobuf:"bytes,5,opt,name=finishTime,proto3" json:"finishTime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateSantriScheduleRequest) Reset() {
	*x = CreateSantriScheduleRequest{}
	mi := &file_santri_schedule_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateSantriScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSantriScheduleRequest) ProtoMessage() {}

func (x *CreateSantriScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSantriScheduleRequest.ProtoReflect.Descriptor instead.
func (*CreateSantriScheduleRequest) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{1}
}

func (x *CreateSantriScheduleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateSantriScheduleRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateSantriScheduleRequest) GetStartPresence() string {
	if x != nil {
		return x.StartPresence
	}
	return ""
}

func (x *CreateSantriScheduleRequest) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *CreateSantriScheduleRequest) GetFinishTime() string {
	if x != nil {
		return x.FinishTime
	}
	return ""
}

type GetSantriScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSantriScheduleRequest) Reset() {
	*x = GetSantriScheduleRequest{}
	mi := &file_santri_schedule_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSantriScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSantriScheduleRequest) ProtoMessage() {}

func (x *GetSantriScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSantriScheduleRequest.ProtoReflect.Descriptor instead.
func (*GetSantriScheduleRequest) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{2}
}

func (x *GetSantriScheduleRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ActiveSantriScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ActiveSantriScheduleRequest) Reset() {
	*x = ActiveSantriScheduleRequest{}
	mi := &file_santri_schedule_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActiveSantriScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActiveSantriScheduleRequest) ProtoMessage() {}

func (x *ActiveSantriScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActiveSantriScheduleRequest.ProtoReflect.Descriptor instead.
func (*ActiveSantriScheduleRequest) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{3}
}

type UpdateSantriScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Schedule      *SantriSchedule        `protobuf:"bytes,1,opt,name=schedule,proto3" json:"schedule,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateSantriScheduleRequest) Reset() {
	*x = UpdateSantriScheduleRequest{}
	mi := &file_santri_schedule_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateSantriScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSantriScheduleRequest) ProtoMessage() {}

func (x *UpdateSantriScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSantriScheduleRequest.ProtoReflect.Descriptor instead.
func (*UpdateSantriScheduleRequest) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateSantriScheduleRequest) GetSchedule() *SantriSchedule {
	if x != nil {
		return x.Schedule
	}
	return nil
}

type DeleteSantriScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteSantriScheduleRequest) Reset() {
	*x = DeleteSantriScheduleRequest{}
	mi := &file_santri_schedule_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteSantriScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSantriScheduleRequest) ProtoMessage() {}

func (x *DeleteSantriScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSantriScheduleRequest.ProtoReflect.Descriptor instead.
func (*DeleteSantriScheduleRequest) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteSantriScheduleRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ListSantriScheduleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSantriScheduleRequest) Reset() {
	*x = ListSantriScheduleRequest{}
	mi := &file_santri_schedule_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSantriScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSantriScheduleRequest) ProtoMessage() {}

func (x *ListSantriScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSantriScheduleRequest.ProtoReflect.Descriptor instead.
func (*ListSantriScheduleRequest) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{6}
}

type ListSantriScheduleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Schedules     []*SantriSchedule      `protobuf:"bytes,1,rep,name=schedules,proto3" json:"schedules,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSantriScheduleResponse) Reset() {
	*x = ListSantriScheduleResponse{}
	mi := &file_santri_schedule_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSantriScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSantriScheduleResponse) ProtoMessage() {}

func (x *ListSantriScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_santri_schedule_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSantriScheduleResponse.ProtoReflect.Descriptor instead.
func (*ListSantriScheduleResponse) Descriptor() ([]byte, []int) {
	return file_santri_schedule_proto_rawDescGZIP(), []int{7}
}

func (x *ListSantriScheduleResponse) GetSchedules() []*SantriSchedule {
	if x != nil {
		return x.Schedules
	}
	return nil
}

var File_santri_schedule_proto protoreflect.FileDescriptor

var file_santri_schedule_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xba, 0x01, 0x0a, 0x0e, 0x53, 0x61, 0x6e, 0x74,
	0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x50, 0x72,
	0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68,
	0x54, 0x69, 0x6d, 0x65, 0x22, 0xb7, 0x01, 0x0a, 0x1b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x2a,
	0x0a, 0x18, 0x47, 0x65, 0x74, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1d, 0x0a, 0x1b, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4a, 0x0a, 0x1b, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x73, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x53, 0x61, 0x6e,
	0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x08, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x22, 0x2d, 0x0a, 0x1b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x1b, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x61, 0x6e, 0x74,
	0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x4b, 0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2d, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x32, 0xc3,
	0x03, 0x0a, 0x15, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x12, 0x1c, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f,
	0x2e, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12,
	0x4d, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1a, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x61, 0x6e, 0x74,
	0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1b, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45,
	0x0a, 0x14, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53,
	0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x3f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x61, 0x6e, 0x74,
	0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x19, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x45, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1c,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x53,
	0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x45, 0x0a,
	0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x61,
	0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x53, 0x61, 0x6e, 0x74, 0x72, 0x69, 0x53, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x64, 0x69, 0x75, 0x62, 0x61, 0x69, 0x64, 0x61, 0x68, 0x2f, 0x72, 0x66,
	0x69, 0x64, 0x2d, 0x73, 0x79, 0x61, 0x66, 0x69, 0x69, 0x79, 0x61, 0x68, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_santri_schedule_proto_rawDescOnce sync.Once
	file_santri_schedule_proto_rawDescData = file_santri_schedule_proto_rawDesc
)

func file_santri_schedule_proto_rawDescGZIP() []byte {
	file_santri_schedule_proto_rawDescOnce.Do(func() {
		file_santri_schedule_proto_rawDescData = protoimpl.X.CompressGZIP(file_santri_schedule_proto_rawDescData)
	})
	return file_santri_schedule_proto_rawDescData
}

var file_santri_schedule_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_santri_schedule_proto_goTypes = []any{
	(*SantriSchedule)(nil),              // 0: SantriSchedule
	(*CreateSantriScheduleRequest)(nil), // 1: CreateSantriScheduleRequest
	(*GetSantriScheduleRequest)(nil),    // 2: GetSantriScheduleRequest
	(*ActiveSantriScheduleRequest)(nil), // 3: ActiveSantriScheduleRequest
	(*UpdateSantriScheduleRequest)(nil), // 4: UpdateSantriScheduleRequest
	(*DeleteSantriScheduleRequest)(nil), // 5: DeleteSantriScheduleRequest
	(*ListSantriScheduleRequest)(nil),   // 6: ListSantriScheduleRequest
	(*ListSantriScheduleResponse)(nil),  // 7: ListSantriScheduleResponse
}
var file_santri_schedule_proto_depIdxs = []int32{
	0, // 0: UpdateSantriScheduleRequest.schedule:type_name -> SantriSchedule
	0, // 1: ListSantriScheduleResponse.schedules:type_name -> SantriSchedule
	1, // 2: SantriScheduleService.CreateSantriSchedule:input_type -> CreateSantriScheduleRequest
	6, // 3: SantriScheduleService.ListSantriSchedule:input_type -> ListSantriScheduleRequest
	3, // 4: SantriScheduleService.ActiveSantriSchedule:input_type -> ActiveSantriScheduleRequest
	2, // 5: SantriScheduleService.GetSantriSchedule:input_type -> GetSantriScheduleRequest
	4, // 6: SantriScheduleService.UpdateSantriSchedule:input_type -> UpdateSantriScheduleRequest
	5, // 7: SantriScheduleService.DeleteSantriSchedule:input_type -> DeleteSantriScheduleRequest
	0, // 8: SantriScheduleService.CreateSantriSchedule:output_type -> SantriSchedule
	7, // 9: SantriScheduleService.ListSantriSchedule:output_type -> ListSantriScheduleResponse
	0, // 10: SantriScheduleService.ActiveSantriSchedule:output_type -> SantriSchedule
	0, // 11: SantriScheduleService.GetSantriSchedule:output_type -> SantriSchedule
	0, // 12: SantriScheduleService.UpdateSantriSchedule:output_type -> SantriSchedule
	0, // 13: SantriScheduleService.DeleteSantriSchedule:output_type -> SantriSchedule
	8, // [8:14] is the sub-list for method output_type
	2, // [2:8] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_santri_schedule_proto_init() }
func file_santri_schedule_proto_init() {
	if File_santri_schedule_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_santri_schedule_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_santri_schedule_proto_goTypes,
		DependencyIndexes: file_santri_schedule_proto_depIdxs,
		MessageInfos:      file_santri_schedule_proto_msgTypes,
	}.Build()
	File_santri_schedule_proto = out.File
	file_santri_schedule_proto_rawDesc = nil
	file_santri_schedule_proto_goTypes = nil
	file_santri_schedule_proto_depIdxs = nil
}
