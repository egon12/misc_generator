package map_graphql_to_proto

import (
	"io"
	"text/template"

	"github.com/vektah/gqlparser/ast"
)

const enumTmpl = `
func MapTo_{{.RetType}}(input {{.ArgType}}) {{.RetType}} {
	return {{.RetType}}(input.String())
}
`

const reverseEnumTmpl = `
func MapToProto_{{.Type}}(input {{.ArgType}}) {{.TypeWithPackage}} {
	value := {{.TypeWithPackage}}_value[string(input)]
	return {{.TypeWithPackage}}(value)
}
`

func generateEnum(input *ast.Definition, output io.Writer) error {
	err := generateEnumOri(input, output)
	if err != nil {
		return err
	}
	return generateReversEnum(input, output)
}

func generateEnumOri(input *ast.Definition, output io.Writer) error {
	config := GetConfig()
	obj := struct {
		ArgType string
		RetType string
	}{}

	obj.ArgType = config.ProtoPackage + "." + removePrefix(input.Name)
	obj.RetType = input.Name

	tmpl := template.New("enumTmpl")
	tmpl, err := tmpl.Parse(enumTmpl)
	if err != nil {
		return err
	}

	return tmpl.Execute(output, obj)
}

func generateReversEnum(input *ast.Definition, output io.Writer) error {
	config := GetConfig()
	obj := struct {
		Type            string
		TypeWithPackage string
		ArgType         string
	}{}

	obj.Type = removePrefix(input.Name)
	obj.TypeWithPackage = config.ProtoPackage + "." + removePrefix(input.Name)
	obj.ArgType = input.Name

	tmpl := template.New("reversEnumTmpl")
	tmpl, err := tmpl.Parse(reverseEnumTmpl)
	if err != nil {
		return err
	}

	return tmpl.Execute(output, obj)
}
