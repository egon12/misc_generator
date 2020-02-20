package graphql_to_proto

import (
	"io"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/ast"
)

const rootResolverTempl = `
func (o *{{.ModuleName}}) {{.FuncName}}(ctx context.Context{{range $a := .Arguments}}, {{$a.Name}} {{$a.Type}}{{end}}) ({{.ReturnType}}, error) {
	return {{.ReturnType}}{}, nil
}
`

const rpcTmpl = `
{{range $f := .RPC}}
message {{.ArgumentType}} {
}{{end}}

service {{.ModuleName}} {
{{range $f := .RPC}}	rpc {{.FuncName}} ({{.ArgumentType}}) returns ({{.ReturnType}});
{{end}}
}

`

func mapRootResolver(queryResolver *ast.Definition, output io.Writer) error {
	tmpl := template.New("rootResolverTmpl")
	tmpl, _ = tmpl.Parse(rpcTmpl)

	type rpc struct {
		FuncName     string
		ArgumentType string
		ReturnType   string
	}

	rr := struct {
		ModuleName string
		RPC        []rpc
	}{"SomeModule", []rpc{}}

	for _, f := range queryResolver.Fields {
		rf := rpc{
			FuncName:     f.Name,
			ArgumentType: strcase.ToCamel(f.Name) + "Request",
		}

		t, _ := mapTypeGraphQltoType(f.Type)
		rf.ReturnType = t

		rr.RPC = append(rr.RPC, rf)
	}

	return tmpl.Execute(output, rr)
}
