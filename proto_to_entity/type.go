package proto_to_entity

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
)

func getEntityType(field *descriptor.FieldDescriptorProto) (string, error) {
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE, descriptor.FieldDescriptorProto_TYPE_ENUM:
		return fixTypeName(field.GetTypeName()), nil
	default:
		return mapToGoEntityType(field.GetType())
	}
}

func fixTypeName(input string) string {
	output := input[1:]

	if output == "google.protobuf.Timestamp" {
		output = "*time.Time"
	} else {
		output = strcase.ToCamel(output)
	}

	return output
}

func mapToGoEntityType(input descriptor.FieldDescriptorProto_Type) (string, error) {

	return mapToPrimitiveGoEntityType(input)
}

func mapToPrimitiveGoEntityType(input descriptor.FieldDescriptorProto_Type) (string, error) {

	switch input {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return reflect.Float64.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return reflect.Float32.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return reflect.Int64.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return reflect.Uint64.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return reflect.Uint32.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return reflect.Int64.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return reflect.Int32.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return reflect.Bool.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return reflect.String.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return "", nil
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return "", nil
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "[]bytes", nil
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return reflect.Uint32.String(), nil
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "", nil
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return "", nil
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return "", nil
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return "", nil
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "", nil
	}

	return "", fmt.Errorf("Cannot find type")
}
