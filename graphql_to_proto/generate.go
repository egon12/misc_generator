package togrpc

import (
	"io"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

func Generate(graphQLContract string, output io.Writer) error {
	return generateProtoFile(graphQLContract, output)
}

func generateProtoFile(graphQLContract string, output io.Writer) error {
	source := &ast.Source{"rootGraphQLContract", graphQLContract, false}

	schema, parseError := parser.ParseSchema(source)
	if parseError != nil {
		return parseError
	}

	var err error
	for _, v := range schema.Definitions {
		switch v.Kind {
		case ast.Object:
			if v.Name == "Query" {
				err = mapRootResolver(v, output)
			} else {
				err = mapResolver(v, output)
			}
		case ast.InputObject:
			err = mapInput(v, output)
		case ast.Enum:
			err = mapEnum(v, output)
		}

		if err != nil {
			return err
		}
	}
	return nil
}
