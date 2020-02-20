package graphql_to_proto

import (
	"fmt"
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
{{range $f := .ArgumentFields}} {{$f}};
{{end}}}
{{end}}

service {{.ModuleName}} {
{{range $f := .RPC}}	rpc {{.FuncName}} ({{.ArgumentType}}) returns ({{.ReturnType}});
{{end}}}
`

func mapRootResolver(queryResolver *ast.Definition, output io.Writer, config Config) error {
	tmpl := template.New("rootResolverTmpl")
	tmpl, _ = tmpl.Parse(rpcTmpl)

	type rpc struct {
		FuncName       string
		ArgumentType   string
		ArgumentFields []string
		ReturnType     string
	}

	rr := struct {
		ModuleName string
		RPC        []rpc
	}{config.getServiceName(), []rpc{}}

	for _, f := range queryResolver.Fields {
		rf := rpc{
			FuncName:       strcase.ToLowerCamel(cleanNameFromPrefix(f.Name, config)),
			ArgumentType:   cleanNameFromPrefix(strcase.ToCamel(f.Name), config) + "Request",
			ArgumentFields: generateArgumentFields(config.RequestAddition, f, config),
		}

		t, _ := mapTypeGraphQltoType(f.Type)
		rf.ReturnType = cleanNameFromPrefix(t, config)

		rr.RPC = append(rr.RPC, rf)
	}

	return tmpl.Execute(output, rr)
}

func generateArgumentFields(requestAddition []string, field *ast.FieldDefinition, config Config) []string {
	lra := len(requestAddition)
	lfa := len(field.Arguments)
	laf := lra + lfa

	result := make([]string, laf, laf)

	for i, ra := range requestAddition {
		result[i] = fmt.Sprintf("%s = %d", ra, i+1)
	}

	for i, af := range field.Arguments {
		t, _ := mapTypeGraphQltoType(af.Type)
		t = cleanNameFromPrefix(t, config)
		result[lra+i] = fmt.Sprintf("%s %s = %d", t, af.Name, lra+i+1)
	}
	return result

}
