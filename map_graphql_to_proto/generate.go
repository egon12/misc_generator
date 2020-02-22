package map_graphql_to_proto

import (
	"io"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

func generateMapper(graphQLContract string, output io.Writer) error {
	s := ast.Source{
		BuiltIn: false,
		Name:    "rootGraphql",
		Input:   graphQLContract,
	}

	sch, parseError := parser.ParseSchema(&s)
	if parseError != nil {
		return parseError
	}

	var err error
	for _, v := range sch.Definitions {
		switch v.Kind {
		case ast.Object:
			if v.Name != "Query" {
				err = generateResolver(v, output)
			}
		case ast.InputObject:
			err = generateInput(v, output)
		}

		if err != nil {
			return err
		}
	}
	return nil
}
