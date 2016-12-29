package genspy

import (
	"fmt"
	"go/ast"
	"go/token"
)

// StructConv converts an interface type into a struct
type StructConv interface {
	IntfToStruct(t *ast.TypeSpec, i *ast.InterfaceType) *ast.TypeSpec
}

// SpyStructConv converts interfaces into spies, i.e., test doubles
type SpyStructConv struct{}

// IntfToStruct mutates the ast.TypeSpec into a struct type with
// public properties for all parameters and return values declared
// in the interface
func (s *SpyStructConv) IntfToStruct(t *ast.TypeSpec, i *ast.InterfaceType) *ast.TypeSpec {
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

	return t
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
				argOffset++
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

				offset++
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
	for idx := range f.Results.List {
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
