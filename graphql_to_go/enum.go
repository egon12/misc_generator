package graphql_to_go

import (
	"io"
	"text/template"

	"github.com/vektah/gqlparser/ast"
)

const enumTmpl = `
type {{.Name}} string

const (
{{range $e := .Enum}}{{$e}} {{$.Name}} = "{{$e}}" 
{{end}})
`

func mapEnum(enum *ast.Definition, output io.Writer) error {
	allEnum := struct {
		Name string
		Enum []string
	}{enum.Name, []string{}}

	for _, en := range enum.EnumValues {
		allEnum.Enum = append(allEnum.Enum, en.Name)
	}

	tmpl := template.New("enumTmpl")
	tmpl, _ = tmpl.Parse(enumTmpl)
	return tmpl.Execute(output, allEnum)
}
