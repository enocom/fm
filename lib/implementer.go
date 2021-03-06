package fm

import (
	"fmt"
	"go/ast"
	"go/token"
)

const (
	recvName = "f"
)

// SpyFuncImplementer creates spy implementations of an interface's functions.
// Meant to be used in conjunction with SpyStructConverter
type SpyFuncImplementer struct{}

// Implement returns a function declaration whose arguments are saved
// as properties and whose return values are properties on a spy struct
func (s *SpyFuncImplementer) Implement(name *ast.Ident, i *ast.InterfaceType) []*ast.FuncDecl {
	var funcDecls []*ast.FuncDecl
	for _, list := range i.Methods.List {
		recv := &ast.FieldList{
			List: []*ast.Field{{
				Names: []*ast.Ident{ast.NewIdent(recvName)},
				Type:  &ast.StarExpr{X: name},
			}},
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
			Body: createBlockStmt(list.Names[0].Name, funcType),
		})
	}
	return funcDecls
}

func createBlockStmt(fname string, f *ast.FuncType) *ast.BlockStmt {
	var list []ast.Stmt

	// x.mu.Lock()
	lockStmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent(recvName),
					Sel: ast.NewIdent("mu"),
				},
				Sel: ast.NewIdent("Lock"),
			},
		},
	}
	list = append(list, lockStmt)

	// defer x.mu.Unlock()
	deferStmt := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent(recvName),
					Sel: ast.NewIdent("mu"),
				},
				Sel: ast.NewIdent("Unlock"),
			},
		},
	}
	list = append(list, deferStmt)

	// add called assignment statement
	calledStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(recvName),
				Sel: ast.NewIdent(fname + "_Called"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ast.NewIdent("true")},
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
				X:   ast.NewIdent(recvName), // for spy
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
