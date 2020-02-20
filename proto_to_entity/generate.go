package proto_to_entity

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
)

func generateFileDescriptorProto(filename string) (*descriptor.FileDescriptorProto, error) {
	fdp := &descriptor.FileDescriptorProto{}

	outputFile := filename + "buf"

	cmd := exec.Command("protoc", "-o"+outputFile, filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fdp, fmt.Errorf("%s\nOutput: %s\n", err.Error(), string(output))
	}

	protobufFile, err := os.Open(outputFile)
	if err != nil {
		return fdp, err
	}

	data, err := ioutil.ReadAll(protobufFile)
	if err != nil {
		return fdp, err
	}

	fdps := &descriptor.FileDescriptorSet{}
	fdps.XXX_Unmarshal(data)

	fdpList := fdps.GetFile()
	if len(fdpList) != 1 {
		return fdp, fmt.Errorf("Want 1 FileDescriptorProto got %d", len(fdpList))
	}

	fdp = fdpList[0]

	return fdp, err
}

func generateStructFromMessage(fdp *descriptor.FileDescriptorProto, output io.Writer) error {

	var b []byte

	for _, msg := range fdp.GetMessageType() {
		structName := strcase.ToCamel(msg.GetName())
		b = []byte(fmt.Sprintf("type %s struct {\n", structName))
		output.Write(b)

		for _, f := range msg.GetField() {
			fieldType, err := getEntityType(f)
			if err != nil {
				return err
			}
			b = []byte(fmt.Sprintf("%s %s\n", strcase.ToCamel(f.GetName()), fieldType))
			output.Write(b)

		}
		output.Write([]byte("}\n"))
	}
	return nil

}

const mapperToEntityFunc = `
func MapToEntity_{{.Name}} (input *proto.{{.Name}}) entity.{{.Name}} {
	return entity.{{.Name}}{
{{ range $f := .Fields }}		{{.Name}}: {{.InputMap}},
{{end}}
	}
}
`

type MapperToEntityField struct {
	Name     string
	InputMap string
}

type MapperToEntity struct {
	Name   string
	Fields []MapperToEntityField
}

func generateMapToEntity(input *descriptor.DescriptorProto, output io.Writer) error {
	var err error
	tmpl := template.New("mapperToEntity")
	tmpl, err = tmpl.Parse(mapperToEntityFunc)
	if err != nil {
		return err
	}

	obj := MapperToEntity{Name: input.GetName()}

	for _, f := range input.GetField() {
		ef := MapperToEntityField{
			Name:     strcase.ToCamel(f.GetName()),
			InputMap: "input." + strcase.ToCamel(f.GetName()),
		}
		obj.Fields = append(obj.Fields, ef)
	}

	return tmpl.Execute(output, obj)
}
