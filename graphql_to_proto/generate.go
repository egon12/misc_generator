package graphql_to_proto

import (
	"io"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type Config struct {
	Prefix          string
	RemovePrefixes  []string
	ServiceName     string
	RequestAddition []string
}

func (c *Config) getServiceName() string {
	if c.ServiceName == "" {
		return "Mod"
	}

	return c.ServiceName
}

func Generate(graphQLContract string, output io.Writer) error {
	return generateProtoFile(graphQLContract, output, Config{})
}

func GenerateWithConfig(graphQLContract string, output io.Writer, config Config) error {
	return generateProtoFile(graphQLContract, output, config)
}

func generateProtoFile(graphQLContract string, output io.Writer, config Config) error {
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
				err = mapRootResolver(v, output, config)
			} else {
				err = mapResolver(v, output, config)
			}
			if err != nil {
				return err
			}
		case ast.InputObject:
			err = mapInput(v, output, config)
			if err != nil {
				return err
			}
		case ast.Enum:
			err = mapEnum(v, output, config)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
