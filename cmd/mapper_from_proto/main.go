package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/egon12/misc_generator/proto_to_entity"
)

func main() {
	mapperType := os.Args[1]
	protoFile := os.Args[2]
	outputFile := os.Args[3]

	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	output.Write([]byte("package mapper\n\n"))
	if mapperType == "e" {
		proto_to_entity.GenerateMapperToEntity(protoFile, output)
	} else if mapperType == "p" {
		proto_to_entity.GenerateMapperToProto(protoFile, output)
	}
	output.Close()

	outputContent, err := exec.Command("goimports", outputFile).CombinedOutput()
	if err != nil {
		log.Fatal(err, string(outputContent))
	}

	output, err = os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	output.Write(outputContent)
	output.Close()
}
