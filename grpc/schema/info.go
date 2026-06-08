package schema

import (
	"sync"

	"github.com/gauas/account-service/dto/response"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

const (
	ServiceName           = "account.v1.AccountService"
	GetUserInfoMethodName = "GetUserInfo"
)

var (
	infoSchemaOnce sync.Once

	getUserInfoRequestDesc  protoreflect.MessageDescriptor
	getUserInfoResponseDesc protoreflect.MessageDescriptor

	requestKeyField protoreflect.FieldDescriptor

	responseKeyField         protoreflect.FieldDescriptor
	responseFullNameField    protoreflect.FieldDescriptor
	responseDobField         protoreflect.FieldDescriptor
	responseGenderField      protoreflect.FieldDescriptor
	responseAvatarURLField   protoreflect.FieldDescriptor
	responseIsOnboardedField protoreflect.FieldDescriptor
	responsePermissionField  protoreflect.FieldDescriptor
)

func NewGetUserInfoRequest() *dynamicpb.Message {
	initInfoSchema()
	return dynamicpb.NewMessage(getUserInfoRequestDesc)
}

func GetUserInfoRequestKey(message *dynamicpb.Message) string {
	initInfoSchema()
	return message.Get(requestKeyField).String()
}

func NewGetUserInfoResponse(profile response.ProfileResponse) *dynamicpb.Message {
	initInfoSchema()

	message := dynamicpb.NewMessage(getUserInfoResponseDesc)
	message.Set(responseKeyField, protoreflect.ValueOfString(profile.Key))
	message.Set(responseFullNameField, protoreflect.ValueOfString(profile.FullName))
	message.Set(responseDobField, protoreflect.ValueOfString(profile.Dob))
	message.Set(responseGenderField, protoreflect.ValueOfString(profile.Gender))
	message.Set(responseAvatarURLField, protoreflect.ValueOfString(profile.AvatarURL))
	message.Set(responseIsOnboardedField, protoreflect.ValueOfBool(profile.IsOnboarded))
	message.Set(responsePermissionField, protoreflect.ValueOfString(profile.Permission))

	return message
}

func initInfoSchema() {
	infoSchemaOnce.Do(func() {
		fileDescriptorProto := &descriptorpb.FileDescriptorProto{
			Syntax:  stringPtr("proto3"),
			Name:    stringPtr("account/v1/account.proto"),
			Package: stringPtr("account.v1"),
			MessageType: []*descriptorpb.DescriptorProto{
				{
					Name: stringPtr("GetUserInfoRequest"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   stringPtr("key"),
							Number: int32Ptr(1),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						},
					},
				},
				{
					Name: stringPtr("GetUserInfoResponse"),
					Field: []*descriptorpb.FieldDescriptorProto{
						stringField("key", 1),
						stringField("full_name", 2),
						stringField("dob", 3),
						stringField("gender", 4),
						stringField("avatar_url", 5),
						boolField("is_onboarded", 6),
						stringField("permission", 7),
					},
				},
			},
			Service: []*descriptorpb.ServiceDescriptorProto{
				{
					Name: stringPtr("AccountService"),
					Method: []*descriptorpb.MethodDescriptorProto{
						{
							Name:       stringPtr(GetUserInfoMethodName),
							InputType:  stringPtr(".account.v1.GetUserInfoRequest"),
							OutputType: stringPtr(".account.v1.GetUserInfoResponse"),
						},
					},
				},
			},
		}

		fileDescriptor, err := protodesc.NewFile(fileDescriptorProto, protoregistry.GlobalFiles)
		if err != nil {
			panic(err)
		}

		if err := protoregistry.GlobalFiles.RegisterFile(fileDescriptor); err != nil {
			panic(err)
		}

		getUserInfoRequestDesc = fileDescriptor.Messages().ByName("GetUserInfoRequest")
		getUserInfoResponseDesc = fileDescriptor.Messages().ByName("GetUserInfoResponse")

		requestKeyField = getUserInfoRequestDesc.Fields().ByName("key")
		responseKeyField = getUserInfoResponseDesc.Fields().ByName("key")
		responseFullNameField = getUserInfoResponseDesc.Fields().ByName("full_name")
		responseDobField = getUserInfoResponseDesc.Fields().ByName("dob")
		responseGenderField = getUserInfoResponseDesc.Fields().ByName("gender")
		responseAvatarURLField = getUserInfoResponseDesc.Fields().ByName("avatar_url")
		responseIsOnboardedField = getUserInfoResponseDesc.Fields().ByName("is_onboarded")
		responsePermissionField = getUserInfoResponseDesc.Fields().ByName("permission")
	})
}

func stringField(name string, number int32) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:   stringPtr(name),
		Number: int32Ptr(number),
		Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
	}
}

func boolField(name string, number int32) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:   stringPtr(name),
		Number: int32Ptr(number),
		Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:   descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(),
	}
}

func stringPtr(value string) *string {
	return &value
}

func int32Ptr(value int32) *int32 {
	return &value
}
