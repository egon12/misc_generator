package map_graphql_to_proto

import (
	"io"
	"text/template"

	"github.com/vektah/gqlparser/ast"
)

func generateSlice(input *ast.SchemaDocument, output io.Writer) error {
	for _, def := range input.Definitions {
		if def.Kind != ast.Object {
			continue
		}

		if def.Name == "Query" {
			continue
		}

		for _, field := range def.Fields {
			if isFieldSlice(field) {
				err := writeSliceMapper(field, output)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func isFieldSlice(field *ast.FieldDefinition) bool {
	return field.Type.Elem != nil
}

const sliceMapperTmpl = `
func MapTo_{{.RetType}}_Slice(input []{{.ArgType}}) []{{.RetType}} {
	n := len(input)
	result := make([]{{.RetType}}, n, n)
	for i, v := range input {
		result[i] = MapTo_{{.RetType}}(v)
	}
	return result
}
`

func writeSliceMapper(field *ast.FieldDefinition, output io.Writer) error {
	config := GetConfig()
	tmpl := template.New("sliceMapperTmpl")
	tmpl, err := tmpl.Parse(sliceMapperTmpl)
	if err != nil {
		return err
	}
	obj := struct {
		RetType string
		ArgType string
	}{}
	obj.RetType = field.Type.Elem.Name()
	obj.ArgType = "*" + config.ProtoPackage + "." + removePrefix(obj.RetType)

	return tmpl.Execute(output, obj)
}
