package main

import (
	"fmt"
	"os"

	. "github.com/egon12/misc_generator/getvar"
)

func main() {
	filename := os.Args[1]
	varname := os.Args[2]

	result, err := GetVar(filename, varname)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
