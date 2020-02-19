package togrpc

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
)

func mapTypeGraphQltoType(input *ast.Type) (string, error) {

	result := ""

	if input.Elem != nil {
		result += result + "repeated "

		child, err := mapTypeGraphQltoType(input.Elem)
		if err != nil {
			return "", fmt.Errorf("Not know type %+v", child)
		}
		return result + child, nil
	}

	if input.NamedType != "" {
		switch input.NamedType {
		case "String":
			return result + "string", nil
		case "Int":
			return result + "int32", nil
		case "Date":
			return result + "google.protobuf.Timestamp", nil
		case "Time":
			return result + "google.protobuf.Timestamp", nil
		}
		return result + input.NamedType, nil
	}

	return "", fmt.Errorf("Not know type %+v", input)
}
