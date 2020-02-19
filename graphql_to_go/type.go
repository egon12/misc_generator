package togoresolver

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
)

func mapTypeGraphQltoType(input *ast.Type) (string, error) {

	result := ""

	if !input.NonNull {
		result = result + "*"
	}

	if input.NamedType != "" {
		switch input.NamedType {
		case "String":
			return result + "string", nil
		case "Int":
			return result + "int32", nil
		case "Date":
			return result + "Date", nil
		case "Time":
			return result + "Time", nil
		}
		return result + input.NamedType, nil
	}

	if input.Elem != nil {
		result += result + "[]"

		child, err := mapTypeGraphQltoType(input.Elem)
		if err != nil {
			return "", fmt.Errorf("Not know type %+v", child)
		}
		return result + child, nil
	}

	return "", fmt.Errorf("Not know type %+v", input)
}
