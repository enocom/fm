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
	pkgs, err := parser.ParseDir(fset, ".", isSrcFile, 0)
	if err != nil {
		fatal(err)
	}

	for _, ast := range pkgs {
		for _, f := range ast.Files {
			decls := findInterfaces(f.Decls)
			f.Decls = decls
			printer.Fprint(os.Stdout, fset, f)
		}
	}
}

func isSrcFile(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

func findInterfaces(ds []ast.Decl) []ast.Decl {
	var decls []ast.Decl
	for _, d := range ds {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			_, ok = typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			decls = append(decls, d)
		}
	}
	return decls
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
