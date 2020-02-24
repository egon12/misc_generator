package graphql_to_go

import (
	"io"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"
)

func mapInput(input *ast.Definition, output io.Writer) error {
	g := gostructStruct{
		Name:   input.Name,
		Fields: []gostructFieldStruct{},
	}

	for _, f := range input.Fields {
		gf := gostructFieldStruct{
			Name:   strcase.ToCamel(f.Name),
			Parent: &g,
		}
		t, _ := mapTypeGraphQltoType(f.Type)
		gf.Type = t
		g.Fields = append(g.Fields, gf)
	}

	tmpl := template.New("gostruct")
	tmpl, err := tmpl.Parse(gostructTmpl)
	if err != nil {
		panic(err)
	}

	return tmpl.Execute(output, g)
}

func removePointer(gotype string) string {
	if strings.HasPrefix(gotype, "*") {
		return strings.Replace(gotype, "*", "", 1)
	}
	return gotype

}
