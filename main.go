package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func main() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", isValidFile, 0)
	if err != nil {
		fatal(err)
	}

	for _, ast := range pkgs {
		for _, f := range ast.Files {
			processFile(fset, f)
		}
	}
}

func isValidFile(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

func processFile(fset *token.FileSet, f *ast.File) {
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

	printer.Fprint(os.Stdout, fset, f)
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
