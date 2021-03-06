package graphql_to_proto

import (
	"io"
	"text/template"

	"github.com/vektah/gqlparser/ast"
)

func mapInput(input *ast.Definition, output io.Writer, config Config) error {
	g := gostructStruct{
		Name:   cleanNameFromPrefix(input.Name, config),
		Fields: []gostructFieldStruct{},
	}

	for _, f := range input.Fields {
		gf := gostructFieldStruct{
			Name:   f.Name,
			Parent: &g,
		}
		t, _ := mapTypeGraphQltoType(f.Type)
		gf.Type = cleanNameFromPrefix(t, config)
		g.Fields = append(g.Fields, gf)
	}

	tmpl := template.New("msg")
	tmpl = tmpl.Funcs(map[string]interface{}{
		"add": func(a, b int) int { return a + b },
	})
	tmpl, err := tmpl.Parse(protoTmpl)

	if err != nil {
		panic(err)
	}

	return tmpl.Execute(output, g)
}
