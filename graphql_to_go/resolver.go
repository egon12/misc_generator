package graphql_to_go

import (
	"io"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"
)

type gostructStruct struct {
	Name   string
	Fields []gostructFieldStruct
}

type gostructFieldStruct struct {
	Name     string
	Type     string
	Parent   *gostructStruct
	FuncName string
}

type rootResolver struct {
	ModuleName string
	FuncName   string
	Arguments  []argument
	ReturnType string
}

type argument struct {
	Name string
	Type string
}

const gostructTmpl = `
type {{.Name}} struct {
{{range $f := .Fields}}{{$f.Name}} {{$f.Type}}
{{end}}}
`

const gostructFuncTmpl = `func (o {{.Parent.Name}}) {{.FuncName}}() {{.Type}} { return o.{{.Name}} }
`

func mapResolver(resolver *ast.Definition, output io.Writer) error {
	g := gostructStruct{
		Name:   resolver.Name,
		Fields: []gostructFieldStruct{},
	}

	for _, f := range resolver.Fields {
		gf := gostructFieldStruct{
			Name:     f.Name,
			FuncName: strcase.ToCamel(f.Name),
			Parent:   &g,
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

	tmpl.Execute(output, g)

	tmpl = template.New("gofunctmpl")
	tmpl, err = tmpl.Parse(gostructFuncTmpl)
	if err != nil {
		panic(err)
	}

	for _, gf := range g.Fields {
		tmpl.Execute(output, gf)
	}

	return nil
}
