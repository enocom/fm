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

	for _, p := range pkgs {
		for _, f := range p.Files {
			interfaceDecls := findInterfaces(f.Decls)
			decls := generateDecls(interfaceDecls)
			// ast.Print(fset, f)

			f.Decls = decls

			f.Name = ast.NewIdent("example_test")
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

func generateDecls(ds []ast.Decl) []ast.Decl {
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

			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			// if we've made it this far, we have an interface to mock

			// start by prefixing the interface's name with "Fake"
			typeSpec.Name = ast.NewIdent("Fake" + typeSpec.Name.Name)

			// convert the interface into a struct
			structType := &ast.StructType{
				Struct:     interfaceType.Interface, // position of the interface keyword
				Fields:     &ast.FieldList{},        // no fields
				Incomplete: interfaceType.Incomplete,
			}
			typeSpec.Type = structType

			// save the decls into the slice
			decls = append(decls, d)

			// add FuncDecls for all the interface's methods
		}
	}
	return decls
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
