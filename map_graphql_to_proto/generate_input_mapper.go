package map_graphql_to_proto

import (
	"io"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"
)

const inputTmpl = `
func MapToProto_{{.Name}}(input {{.ArgType}}) *{{.RetType}} {
	return &{{.RetType}} {
	{{range $f := .Fields}}	{{ $f.Output }}: {{ $f.Input }},
	{{end}}}
}`

type fieldMapper struct {
	Input  string
	Output string
}

func generateInput(input *ast.Definition, output io.Writer) error {
	// sory for global variable
	config := GetConfig()

	obj := struct {
		Name    string
		ArgType string
		RetType string
		Fields  []fieldMapper
	}{}

	obj.ArgType = input.Name
	obj.RetType = config.ProtoPackage + "." + removePrefix(input.Name)
	obj.Name = removePrefix(input.Name)

	for _, f := range input.Fields {
		fm := fieldMapper{
			Input:  getToProtoInputFunction(f),
			Output: strcase.ToCamel(f.Name),
		}
		obj.Fields = append(obj.Fields, fm)
	}

	tmpl := template.New("inputTmpl")
	tmpl, err := tmpl.Parse(inputTmpl)
	if err != nil {
		return err
	}

	return tmpl.Execute(output, obj)
}
