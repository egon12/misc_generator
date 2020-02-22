package proto_to_entity

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/strcase"
)

func GenerateEntity(protoFileName string, output io.Writer) error {
	fdp, err := generateFileDescriptorProto(protoFileName)
	if err != nil {
		return err
	}
	return generateStructFromMessage(fdp, output)
}

func GenerateMapperToEntity(protoFileName string, output io.Writer) error {
	fdp, err := generateFileDescriptorProto(protoFileName)
	if err != nil {
		return err
	}
	return generateMapToEntitis(fdp, output)
}

func GenerateMapperToProto(protoFileName string, output io.Writer) error {
	fdp, err := generateFileDescriptorProto(protoFileName)
	if err != nil {
		return err
	}
	return generateMapToProtos(fdp, output)
}

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

	packageName := fdp.GetPackage()

	for _, msg := range fdp.GetMessageType() {
		err := generateStructFromSingleMessage(msg, output, packageName)
		if err != nil {
			return err
		}
		/*
			structName := strcase.ToCamel(msg.GetName())
			b = []byte(fmt.Sprintf("type %s struct {\n", structName))
			output.Write(b)

			for _, f := range msg.GetField() {

				fieldType, err := getEntityType(f, packageName)
				if err != nil {
					return err
				}
				b = []byte(fmt.Sprintf("%s %s\n", strcase.ToCamel(f.GetName()), fieldType))
				output.Write(b)

			}
			output.Write([]byte("}\n"))
		*/
	}
	return nil

}

func generateStructFromSingleMessage(input *descriptor.DescriptorProto, output io.Writer, packageName string) error {
	var b []byte

	structName := strcase.ToCamel(input.GetName())
	b = []byte(fmt.Sprintf("type %s struct {\n", structName))
	output.Write(b)

	for _, f := range input.GetField() {

		fieldType, err := getEntityType(f, packageName)
		if err != nil {
			return err
		}
		b = []byte(fmt.Sprintf("%s %s\n", strcase.ToCamel(f.GetName()), fieldType))
		output.Write(b)

	}
	output.Write([]byte("}\n"))
	return nil
}
