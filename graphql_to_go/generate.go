package graphql_to_go

import (
	"fmt"
	"io"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

func Generate(graphQLContract string, output io.Writer) error {
	err := generateResolver(graphQLContract, output)
	if err != nil {
		return err
	}
	return generateRootResolver(graphQLContract, output)
}

func GenerateResolver(graphQLContract string, output io.Writer) error {
	return generateResolver(graphQLContract, output)
}

func GenerateRootResolver(graphQLContract string, output io.Writer) error {
	return generateRootResolver(graphQLContract, output)
}

func generateResolver(graphQLContract string, output io.Writer) error {
	s := ast.Source{
		BuiltIn: true,
		Name:    "rootGraphql",
		Input:   graphQLContract,
	}

	var err error
	sch, parseError := parser.ParseSchema(&s)
	if parseError != nil {
		return parseError
	}

	for _, v := range sch.Definitions {
		switch v.Kind {
		case ast.Object:
			if v.Name != "Query" {
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

func generateRootResolver(graphQLContract string, output io.Writer) error {
	s := ast.Source{
		BuiltIn: true,
		Name:    "rootGraphql",
		Input:   graphQLContract,
	}

	sch, parseError := parser.ParseSchema(&s)
	if parseError != nil {
		return parseError
	}

	for _, v := range sch.Definitions {
		if v.Kind == ast.Object && v.Name == "Query" {
			return mapRootResolver(v, output)
		}
	}
	return fmt.Errorf("Cannot find type Query")
}
