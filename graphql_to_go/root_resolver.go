package graphql_to_go

import (
	"io"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"
)

const rootResolverTempl = `
func (o *{{.ModuleName}}) {{.FuncName}}(
	ctx context.Context,
	{{if .Arguments}}args struct {
		{{range $a := .Arguments}}{{$a.Name}} {{$a.Type}}
		{{end}}},
	{{end}}
) ({{.ReturnType}}, error) {
	return {{.ReturnType}}{}, nil
}
`

func mapRootResolver(queryResolver *ast.Definition, output io.Writer) error {
	tmpl := template.New("rootResolverTmpl")
	tmpl, _ = tmpl.Parse(rootResolverTempl)

	for _, f := range queryResolver.Fields {
		rr := rootResolver{
			ModuleName: "SomeModuleName",
			FuncName:   strcase.ToCamel(f.Name),
			Arguments:  []argument{},
		}

		for _, a := range f.Arguments {
			t, _ := mapTypeGraphQltoType(a.Type)
			rr.Arguments = append(rr.Arguments, argument{
				Name: strcase.ToCamel(a.Name),
				Type: t,
			})
		}

		t, _ := mapTypeGraphQltoType(f.Type)
		rr.ReturnType = t
		tmpl.Execute(output, rr)
	}

	return nil
}
