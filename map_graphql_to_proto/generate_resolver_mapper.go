package map_graphql_to_proto

import (
	"io"
	"text/template"

	"github.com/vektah/gqlparser/ast"
)

const resolverTmpl = `
func MapTo_{{.Name}}(input {{.ArgType}}) {{.RetType}} {
	return {{.RetType}} {
	{{range $f := .Fields}}	{{ $f.Output }}: {{ $f.Input }},
	{{end}}}
}`

func generateResolver(input *ast.Definition, output io.Writer) error {
	config := GetConfig()
	obj := struct {
		Name    string
		ArgType string
		RetType string
		Fields  []fieldMapper
	}{}

	obj.ArgType = "*" + config.ProtoPackage + "." + removePrefix(input.Name)
	obj.RetType = input.Name
	obj.Name = input.Name

	for _, f := range input.Fields {
		fm := fieldMapper{
			Input:  getToGraphQLInputFunction(f),
			Output: f.Name,
		}
		obj.Fields = append(obj.Fields, fm)
	}

	tmpl := template.New("resolveTmpl")
	tmpl, err := tmpl.Parse(resolverTmpl)
	if err != nil {
		return err
	}

	return tmpl.Execute(output, obj)
}
