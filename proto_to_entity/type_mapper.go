package proto_to_entity

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
)

func getMapperToEntityField(field *descriptor.FieldDescriptorProto, packageName string) (string, error) {

	if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED && field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return "MapToEntity_" + fixTypeName(field.GetTypeName(), packageName) + "Slice(input.Get" + strcase.ToCamel(field.GetName()) + "())", nil
	}

	switch field.GetType() {

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "input.Get" + strcase.ToCamel(field.GetName()) + "()", nil

	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32:
		return "int(input.Get" + strcase.ToCamel(field.GetName()) + "())", nil

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "entity." + fixTypeName(field.GetTypeName(), packageName) + "(input.Get" + strcase.ToCamel(field.GetName()) + "())", nil
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		fixType := fixTypeName(field.GetTypeName(), packageName)
		if fixType == "*time.Time" {
			return "FromTimestamp(input.Get" + strcase.ToCamel(field.GetName()) + "())", nil
		}
		return "MapToEntity_" + fixTypeName(field.GetTypeName(), packageName) + "(input.Get" + strcase.ToCamel(field.GetName()) + "())", nil
	default:
		return "", fmt.Errorf("MapToEntity for %s not Yet implemented", field.GetName())
	}
}

func getMapperToProtoField(field *descriptor.FieldDescriptorProto, packageName string) (string, error) {

	if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED && field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return "MapToProto_" + fixTypeName(field.GetTypeName(), packageName) + "Slice(input." + strcase.ToCamel(field.GetName()) + ")", nil
	}

	switch field.GetType() {

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "input." + strcase.ToCamel(field.GetName()), nil

	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return "int64(input." + strcase.ToCamel(field.GetName()) + ")", nil

	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32:
		return "int32(input." + strcase.ToCamel(field.GetName()) + ")", nil

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "proto." + fixTypeName(field.GetTypeName(), packageName) + "(input." + strcase.ToCamel(field.GetName()) + ")", nil
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		fixType := fixTypeName(field.GetTypeName(), packageName)
		if fixType == "*time.Time" {
			return "ToTimestamp(input." + strcase.ToCamel(field.GetName()) + ")", nil
		}
		return "MapToProto_" + fixTypeName(field.GetTypeName(), packageName) + "(input." + strcase.ToCamel(field.GetName()) + ")", nil
	default:
		return "", fmt.Errorf("MapToProto for %s not Yet implemented", field.GetName())
	}

}
