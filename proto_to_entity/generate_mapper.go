package proto_to_entity

import (
	"io"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
)

const mapperToEntityFunc = `
func MapToEntity_{{.Name}} (input *proto.{{.Name}}) entity.{{.Name}} {
	return entity.{{.Name}}{
{{ range $f := .Fields }}		{{.Name}}: {{.InputMap}},
{{end}}
	}
}
`

const mapperToProtoFunc = `
func MapToProto_{{.Name}} (input entity.{{.Name}}) *proto.{{.Name}} {
	return &proto.{{.Name}}{
{{ range $f := .Fields }}		{{.Name}}: {{.InputMap}},
{{end}}
	}
}
`

type mapperToEntityField struct {
	Name     string
	InputMap string
}

type mapperToEntity struct {
	Name   string
	Fields []mapperToEntityField
}

func generateMapToEntitis(input *descriptor.FileDescriptorProto, output io.Writer) (err error) {
	packageName := input.GetPackage()
	for _, i := range input.GetMessageType() {
		err = generateMapToEntity(i, output, packageName)
		if err != nil {
			return
		}
	}

	err = generateMapToSliceEntities(input, output)
	return
}

func generateMapToEntity(input *descriptor.DescriptorProto, output io.Writer, packageName string) error {
	var err error
	tmpl := template.New("mapperToEntity")
	tmpl, err = tmpl.Parse(mapperToEntityFunc)
	if err != nil {
		return err
	}

	obj := mapperToEntity{Name: input.GetName()}

	for _, f := range input.GetField() {
		inputMap, err := getMapperToEntityField(f, packageName)
		if err != nil {
			return err
		}
		ef := mapperToEntityField{
			Name:     strcase.ToCamel(f.GetName()),
			InputMap: inputMap,
		}
		obj.Fields = append(obj.Fields, ef)
	}

	return tmpl.Execute(output, obj)
}

func generateMapToProtos(input *descriptor.FileDescriptorProto, output io.Writer) (err error) {
	packageName := input.GetPackage()
	for _, i := range input.GetMessageType() {
		err = generateMapToProto(i, output, packageName)
		if err != nil {
			return
		}
	}
	err = generateMapToSliceProtos(input, output)
	return
}

func generateMapToProto(input *descriptor.DescriptorProto, output io.Writer, packageName string) error {
	var err error
	tmpl := template.New("mapperToProto")
	tmpl, err = tmpl.Parse(mapperToProtoFunc)
	if err != nil {
		return err
	}

	obj := mapperToEntity{Name: input.GetName()}

	for _, f := range input.GetField() {
		inputMap, err := getMapperToProtoField(f, packageName)
		if err != nil {
			return err
		}
		ef := mapperToEntityField{
			Name:     strcase.ToCamel(f.GetName()),
			InputMap: inputMap,
		}
		obj.Fields = append(obj.Fields, ef)
	}

	return tmpl.Execute(output, obj)
}
