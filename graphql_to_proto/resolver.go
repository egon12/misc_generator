package graphql_to_proto

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
	Arguments  []Argument
	ReturnType string
}

type Argument struct {
	Name string
	Type string
}

const gostructTmpl = `
type {{.Name}} struct {
{{range $f := .Fields}}{{$f.Name}} {{$f.Type}}
{{end}}}
`

const gostructFuncTmpl = `
func (o *{{.Parent.Name}}) {{.FuncName}}() {{.Type}} {
	return o.{{.Name}}
}
`

const protoTmpl = `
message {{.Name}} {
{{range $i, $f := .Fields}}	{{$f.Type}} {{$f.Name}} = {{ add $i 1}};
{{end}}}
`

func mapResolver(resolver *ast.Definition, output io.Writer, config Config) error {
	g := gostructStruct{
		Name:   cleanNameFromPrefix(resolver.Name, config),
		Fields: []gostructFieldStruct{},
	}

	for _, f := range resolver.Fields {
		gf := gostructFieldStruct{
			Name:     f.Name,
			FuncName: cleanNameFromPrefix(strcase.ToCamel(f.Name), config),
			Parent:   &g,
		}
		t, _ := mapTypeGraphQltoType(f.Type)
		gf.Type = cleanNameFromPrefix(t, config)
		g.Fields = append(g.Fields, gf)
	}

	tmpl := template.New("prototmpl")
	tmpl = tmpl.Funcs(map[string]interface{}{
		"add": func(a, b int) int { return a + b },
	})
	tmpl, err := tmpl.Parse(protoTmpl)
	if err != nil {
		panic(err)
	}

	tmpl.Execute(output, g)

	return nil
}
