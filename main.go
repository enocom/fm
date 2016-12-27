package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	if err := filepath.Walk(dir, processFile); err != nil {
		fatal(err)
	}
}

func isValidFile(info os.FileInfo) bool {
	// temporary
	if info.Name() == "main.go" {
		return false
	}

	if info.IsDir() {
		return false
	}

	return strings.HasSuffix(info.Name(), ".go") &&
		!strings.HasSuffix(info.Name(), "_test.go")
}

func processFile(path string, info os.FileInfo, err error) error {
	if !isValidFile(info) {
		return nil
	}

	// create ast from file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, info.Name(), nil, 0)
	if err != nil {
		fatal(err)
	}

	var decls []ast.Decl
	for _, d := range f.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)

			if !ok {
				continue
			}

			if _, ok := typeSpec.Type.(*ast.InterfaceType); !ok {
				continue
			}

			decls = append(decls, d)
		}
	}
	f.Decls = decls

	return printer.Fprint(os.Stdout, fset, f)
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
