package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const (
	fakePrefix   = "Fake"
	recvName     = "f"
	inputSuffix  = "_Input"
	outputSuffix = "_Output"
	argPrefix    = "Arg"
	retPrefix    = "Ret"
)

// NewGenerator creates a generator which will find all interfaces
// within `workingDir`, create spy implementations of those interfaces,
// and then write the result out to the file named by `fileDst`
func NewGenerator(workingDir, fileDst string) *generator {
	return &generator{wd: workingDir, dst: fileDst}
}

type generator struct {
	wd  string
	dst string
}

func (g *generator) Run() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, g.wd, isSrcFile, 0)
	if err != nil {
		fatal(err)
	}

	for pname, p := range pkgs {
		spyFile, err := os.Create(g.dst)
		if err != nil {
			fatal(err)
		}

		var decls []ast.Decl
		for _, f := range p.Files {
			decls = append(decls, generateSpyDecls(f.Decls)...)
		}

		astFile := &ast.File{
			Name:  ast.NewIdent(pname + "_test"),
			Decls: decls,
		}

		format.Node(spyFile, fset, astFile)
	}
}

func isSrcFile(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

// generateSpies transforms all the interfaces in the list of declarations
// into spies in the form of structs with implemented functions
func generateSpyDecls(ds []ast.Decl) []ast.Decl {
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
	t.Name = ast.NewIdent(fakePrefix + t.Name.Name)

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
			inputStruct := buildStruct(methodName+inputSuffix, argPrefix, funcType.Params.List)
			list = append(list, inputStruct)
		}

		// add Output struct with result values
		if len(funcType.Results.List) > 0 {
			outputStruct := buildStruct(methodName+outputSuffix, retPrefix, funcType.Results.List)
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
					Names: []*ast.Ident{ast.NewIdent(recvName)},
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
	var list []ast.Stmt

	// add called assignment statement
	calledStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(recvName),
				Sel: ast.NewIdent(fname + "_Called"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			ast.NewIdent("true"),
		},
	}
	list = append(list, calledStmt)

	// add assignment for each param
	offset := 0
	for idx, field := range f.Params.List {
		// make assignments for multiple fields with same type
		if len(field.Names) > 1 {
			for _, name := range field.Names {
				assignStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(recvName),
								Sel: ast.NewIdent(fname + inputSuffix),
							},
							Sel: ast.NewIdent(fmt.Sprintf("%s%d", argPrefix, idx+offset)),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{name},
				}
				list = append(list, assignStmt)

				offset += 1
			}
		} else {
			assignStmt := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent(recvName),
							Sel: ast.NewIdent(fname + inputSuffix),
						},
						Sel: ast.NewIdent(fmt.Sprintf("%s%d", argPrefix, idx+offset)),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{field.Names[0]},
			}
			list = append(list, assignStmt)
		}
	}

	// add return statement if there are values to return
	var results []ast.Expr
	for idx, _ := range f.Results.List {
		results = append(results, &ast.SelectorExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(recvName), // for fake
				Sel: ast.NewIdent(fname + outputSuffix),
			},
			Sel: ast.NewIdent(fmt.Sprintf("%s%d", retPrefix, idx)),
		})
	}
	if len(results) > 0 {
		list = append(list, &ast.ReturnStmt{Results: results})
	}

	return &ast.BlockStmt{List: list}
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
