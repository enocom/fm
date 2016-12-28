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

	// TODO: ensure multiple files output into single file
	for name, p := range pkgs {
		for _, f := range p.Files {
			// ast.Print(fset, f)
			f.Decls = generateSpies(f.Decls)
			f.Name = ast.NewIdent(name + "_test")

			// TODO: Write to file instead of standard out
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

			mutatateToStruct(typeSpec, interfaceType)
			funcDecls := createSpyFuncs(typeSpec, interfaceType)

			decls = append(decls, genDecl)
			for _, fd := range funcDecls {
				decls = append(decls, fd)
			}
		}
	}
	return decls
}

// mutatateToStruct mutates the underlying interface type into a struct type
// and adds implentations of the interface's methods
func mutatateToStruct(t *ast.TypeSpec, i *ast.InterfaceType) {
	t.Name = ast.NewIdent("Fake" + t.Name.Name)

	var list []*ast.Field
	for _, field := range i.Methods.List {
		funcType, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue // TODO: when would this happen?
		}

		methodName := field.Names[0].Name
		wasCalled := &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(methodName + "_Called")},
			Type:  ast.NewIdent("bool"),
		}
		list = append(list, wasCalled)

		// add Input struct with params
		if len(funcType.Params.List) > 0 {
			input := buildInputStruct("Arg", methodName+"_Input", funcType.Params.List)
			list = append(list, input)
		}

		// add Output struct with return values
		if len(funcType.Results.List) > 0 {
			// for idx, result := range funcType.Results.List {
			// }
		}
	}

	t.Type = &ast.StructType{
		Fields:     &ast.FieldList{List: list},
		Incomplete: i.Incomplete,
	}
}

// buildInputStruct writes a struct type whose fields
// reflect the various input arguments defined in the interface
func buildInputStruct(prefix, fieldname string, list []*ast.Field) *ast.Field {
	var fields []*ast.Field
	var argOffset int
	var argName string

	for idx, param := range list {
		argName = fmt.Sprintf("%s%d", prefix, idx+argOffset)
		// if we have multiple args of same type,
		// add fields for each
		if len(param.Names) > 1 {
			for range param.Names {
				argName = fmt.Sprintf("%s%d", prefix, idx+argOffset)
				fields = append(fields, &ast.Field{
					Names: []*ast.Ident{ast.NewIdent(argName)},
					Type:  param.Type,
				})
				argOffset += 1
			}
		} else {
			fields = append(fields, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(argName)},
				Type:  param.Type,
			})
		}
	}
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(fieldname)},
		Type: &ast.StructType{
			Fields: &ast.FieldList{List: fields},
		},
	}
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
			// TODO: When will this happen?
			continue
		}

		blockStmt := &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						// TODO: stop hard coding these values
						&ast.BasicLit{Kind: token.INT, Value: "0"},
						ast.NewIdent("nil"),
					},
				},
			},
		}

		funcDecls = append(funcDecls, &ast.FuncDecl{
			Recv: recv,
			Name: list.Names[0],
			Type: funcType,
			Body: blockStmt,
		})
	}
	return funcDecls
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
