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
			// time to create a spy
			createSpyStruct(typeSpec, interfaceType)

			// add implentations of the interface's methods
			funcDecls := createSpyFuncs(typeSpec, interfaceType)

			decls = append(decls, genDecl)
			for _, fd := range funcDecls {
				decls = append(decls, fd)
			}
		}
	}
	return decls
}

// createSpy mutates the underlying interface type into a struct type
// and adds implentations of the interface's methods
func createSpyStruct(t *ast.TypeSpec, i *ast.InterfaceType) {
	// start by prefixing the interface's name with "Fake"
	t.Name = ast.NewIdent("Fake" + t.Name.Name)

	// convert the interface into a struct
	structType := &ast.StructType{
		Struct:     i.Interface,      // position of the interface keyword
		Fields:     &ast.FieldList{}, // no fields
		Incomplete: i.Incomplete,
	}
	t.Type = structType
}

// createSpyFuncs creates spy functions which implement the methods of
// the interface type
func createSpyFuncs(t *ast.TypeSpec, i *ast.InterfaceType) []*ast.FuncDecl {
	var funcDecls []*ast.FuncDecl
	for _, list := range i.Methods.List {
		recv := &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{ast.NewIdent("f")},
					Type:  &ast.StarExpr{X: t.Name},
				},
			},
		}

		funcType, ok := list.Type.(*ast.FuncType)
		if !ok {
			// the type cast will fail on embedded types
			// ignoring for now
			continue
		}

		blockStmt := &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{Kind: token.INT, Value: "0"},
						ast.NewIdent("nil"),
					},
				},
			},
		}

		funcDecls = append(funcDecls, &ast.FuncDecl{
			Recv: recv,          // *FieldList
			Name: list.Names[0], // *Ident
			Type: funcType,      // *FuncType
			Body: blockStmt,     // *BlockStmt
		})
	}
	return funcDecls
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
