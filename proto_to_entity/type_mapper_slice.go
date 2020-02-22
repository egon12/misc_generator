package proto_to_entity

import (
	"io"
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

const sliceMapperToEntityFunc = `
func MapToEntity_{{.Name}}Slice(input []*proto.{{.Name}}) (result []entity.{{.Name}}) {
	result = make([]entity.{{.Name}}, len(input), len(input))
	for i, v := range input {
		result[i] = MapToEntity_{{.Name}}(v)
	}
	return
}
`

const sliceMapperToProtoFunc = `
func MapToProto_{{.Name}}Slice(input []entity.{{.Name}}) (result []*proto.{{.Name}}) {
	result = make([]*proto.{{.Name}}, len(input), len(input))
	for i, v := range input {
		result[i] = MapToProto_{{.Name}}(v)
	}
	return
}
`

func generateMapToSliceEntities(input *descriptor.FileDescriptorProto, output io.Writer) (err error) {
	tmpl := template.New("sliceMapperToEntityFunc")
	tmpl, err = tmpl.Parse(sliceMapperToEntityFunc)
	if err != nil {
		return
	}

	input.GetPackage()
	for _, m := range input.GetMessageType() {
		for _, f := range m.GetField() {
			if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				var t string
				t, err = getEntityType(f, input.GetPackage())
				if err != nil {
					return
				}
				t = strings.Replace(t, "[]", "", 1)
				err = tmpl.Execute(output, struct{ Name string }{t})
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func generateMapToSliceProtos(input *descriptor.FileDescriptorProto, output io.Writer) (err error) {
	tmpl := template.New("sliceMapperToProtoFunc")
	tmpl, err = tmpl.Parse(sliceMapperToProtoFunc)
	if err != nil {
		return
	}

	input.GetPackage()
	for _, m := range input.GetMessageType() {
		for _, f := range m.GetField() {
			if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				var t string
				t, err = getEntityType(f, input.GetPackage())
				if err != nil {
					return
				}
				t = strings.Replace(t, "[]", "", 1)
				err = tmpl.Execute(output, struct{ Name string }{t})
				if err != nil {
					return
				}
			}
		}
	}

	return
}
