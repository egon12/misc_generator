package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/egon12/misc_generator/proto_to_entity"
)

func main() {
	protoFile := os.Args[1]
	outputFile := os.Args[2]

	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	output.Write([]byte("package entity\n\n"))
	proto_to_entity.GenerateEntity(protoFile, output)
	output.Close()

	outputContent, err := exec.Command("goimports", outputFile).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	output, err = os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	output.Write(outputContent)
	output.Close()
}
