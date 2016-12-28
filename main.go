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

			mutateToStruct(typeSpec, interfaceType)
			funcDecls := createSpyFuncs(typeSpec, interfaceType)

			decls = append(decls, genDecl)
			for _, fd := range funcDecls {
				decls = append(decls, fd)
			}
		}
	}
	return decls
}

// mutateToStruct mutates the underlying interface type into a struct type
// and adds implentations of the interface's methods
func mutateToStruct(t *ast.TypeSpec, i *ast.InterfaceType) {
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

		// add Input struct with arguments
		if len(funcType.Params.List) > 0 {
			inputStruct := buildStruct(methodName+"_Input", "Arg", funcType.Params.List)
			list = append(list, inputStruct)
		}

		// add Output struct with result values
		if len(funcType.Results.List) > 0 {
			outputStruct := buildStruct(methodName+"_Output", "Ret", funcType.Results.List)
			list = append(list, outputStruct)
		}
	}

	t.Type = &ast.StructType{
		Fields:     &ast.FieldList{List: list},
		Incomplete: i.Incomplete,
	}
}

// buildInputStruct writes a struct type whose fields
// reflect the various input arguments defined in the interface
func buildStruct(fieldname, prefix string, list []*ast.Field) *ast.Field {
	var fields []*ast.Field
	var argOffset int
	var argName string

	for idx, param := range list {
		// if we have multiple arguments of the same type,
		// add fields for each argument
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
			argName = fmt.Sprintf("%s%d", prefix, idx+argOffset)
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
					Names: []*ast.Ident{ast.NewIdent("f")}, // "f" for fake
					Type:  &ast.StarExpr{X: t.Name},
				},
			},
		}

		funcType, ok := list.Type.(*ast.FuncType)
		if !ok {
			// TODO: When will this happen?
			continue
		}

		funcDecls = append(funcDecls, &ast.FuncDecl{
			Recv: recv,
			Name: list.Names[0],
			Type: funcType,
			Body: createBlockStmt(t, list.Names[0].Name, funcType),
		})
	}
	return funcDecls
}

func createBlockStmt(t *ast.TypeSpec, fname string, f *ast.FuncType) *ast.BlockStmt {
	// TODO: for each Param, save it off in the corresponding Input field
	// TODO: return set values for specified in Output
	calledStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("f"),
				Sel: ast.NewIdent(fname + "_Called"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			ast.NewIdent("true"),
		},
	}

	returnStmt := &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.BasicLit{Kind: token.INT, Value: "0"},
			ast.NewIdent("nil"),
		},
	}

	return &ast.BlockStmt{
		List: []ast.Stmt{calledStmt, returnStmt},
	}
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
