package getvar

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// GetVar get variable value from file by specify fileName and varName
func GetVar(fileName, varName string) (string, error) {
	fs := token.NewFileSet()

	f, err := parser.ParseFile(fs, fileName, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	for _, d := range f.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if !(genDecl.Tok == token.VAR || genDecl.Tok == token.CONST) {
			continue
		}

		for _, s := range genDecl.Specs {
			valueSpec, ok := s.(*ast.ValueSpec)
			if !ok {
				continue
			}

			if valueSpec.Names[0].Name != varName {
				continue
			}

			basicLit, ok := valueSpec.Values[0].(*ast.BasicLit)
			if !ok {
				continue
			}

			return basicLit.Value, nil
		}
	}

	return "", fmt.Errorf("Cannot find %s in %s", varName, fileName)
}
