package map_graphql_to_proto

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"
)

func getToGraphQLInputFunction(input *ast.FieldDefinition) string {

	if input.Type.Elem != nil {
		return fmt.Sprintf("MapTo_%s_Slice(input.Get%s())",
			input.Type.Elem.Name(),
			strcase.ToCamel(removePrefix(input.Name)),
		)
	}

	switch input.Type.NamedType {
	case
		"String",
		"Bool":
		if input.Type.NonNull {
			return getGetFunction(input)
		} else {
			return "&input." + strcase.ToCamel(input.Name)
		}
	case "Int":
		return fmt.Sprintf("int32(%s)", getGetFunction(input))
	case "Time":
		if input.Type.NonNull {
			return fmt.Sprintf("FromTimestampProtoNotNil(%s)", getGetFunction(input))
		} else {
			return fmt.Sprintf("FromTimestampProto(%s)", getGetFunction(input))
		}
	default:
		return fmt.Sprintf("MapTo_%s(input.Get%s())",
			input.Type.NamedType,
			strcase.ToCamel(removePrefix(input.Name)),
		)
	}
}

func getGetFunction(input *ast.FieldDefinition) string {
	return "input.Get" + strcase.ToCamel(input.Name) + "()"
}

func getToProtoInputFunction(input *ast.FieldDefinition) string {
	name := strcase.ToCamel(input.Name)
	inputType := input.Type.NamedType

	ptrPrefix := ""
	if !input.Type.NonNull {
		ptrPrefix = "*"
	}

	switch inputType {
	case "Date":
		return "input." + name + ".TimestampProto()"

	case "Int", "String", "Bool":
		return ptrPrefix + "input." + name

	default:
		return fmt.Sprintf("MapToProto_%s(%sinput.%s)",
			removePrefix(inputType),
			ptrPrefix,
			name,
		)
	}
}
