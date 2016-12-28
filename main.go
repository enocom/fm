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
			// ast.Print(fset, f)
			f.Decls = generateSpies(f.Decls)
			f.Name = ast.NewIdent("example_test")

			printer.Fprint(os.Stdout, fset, f)
		}
	}
}

func isSrcFile(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

// generateSpies transforms all the interfaces in the list of declarations
// into spies in the form of structs with implemented functions
func generateSpies(ds []ast.Decl) []ast.Decl {
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

			// found an interface
			// time to create a mock
			createSpy(typeSpec, interfaceType)

			decls = append(decls, genDecl)
		}
	}
	return decls
}

// createSpy mutates the underlying interface type into a struct type
// and adds implentations of the interface's methods
func createSpy(typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) {
	// start by prefixing the interface's name with "Fake"
	typeSpec.Name = ast.NewIdent("Fake" + typeSpec.Name.Name)

	// convert the interface into a struct
	structType := &ast.StructType{
		Struct:     interfaceType.Interface, // position of the interface keyword
		Fields:     &ast.FieldList{},        // no fields
		Incomplete: interfaceType.Incomplete,
	}
	typeSpec.Type = structType
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
