package fm_test

import (
	"go/ast"
	"testing"

	fm "github.com/enocom/fm/lib"
)

// TestGenerateReturnsSliceOfSpyDecls ensures the generator produces
// two declarations for a single interface with a single method:
// 1) a struct with fields to store the result of a function call, and
// 2) a spy implementation of the interface's single method.
func TestGenerateReturnsSliceOfSpyDecls(t *testing.T) {
	gen := &fm.SpyGenerator{
		Converter:   &fm.SpyStructConverter{},
		Implementer: &fm.SpyFuncImplementer{},
	}
	interfaceDecls := buildInterfaceAST()
	spyDecls := gen.Generate(interfaceDecls)

	want := 2
	got := len(spyDecls)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGenerateSkipsAnythingButInterfaceTypes(t *testing.T) {
	gen := &fm.SpyGenerator{
		Converter:   &fm.SpyStructConverter{},
		Implementer: &fm.SpyFuncImplementer{},
	}
	decls := buildASTWithoutInterfaces()
	spyDecls := gen.Generate(decls)

	want := 0
	got := len(spyDecls)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGenerateReturnsEmptySliceForNoInput(t *testing.T) {
	gen := &fm.SpyGenerator{}

	result := gen.Generate(make([]ast.Decl, 0))

	want := 0
	got := len(result)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// buildTestAST generates as AST for the following interface:
//
// type Tester interface {
//     Test()
// }
func buildInterfaceAST() []ast.Decl {
	fields := make([]*ast.Field, 0)
	fields = append(fields, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("Test")},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: nil},
			Results: &ast.FieldList{List: nil},
		},
	})

	decls := make([]ast.Decl, 0)
	decls = append(decls, &ast.GenDecl{
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Tester"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	})

	return decls
}

// buildASTWithoutInterfaces generates as AST of the following code:
//
// type Tester struct {}
func buildASTWithoutInterfaces() []ast.Decl {
	decls := make([]ast.Decl, 0)
	decls = append(decls, &ast.GenDecl{
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Tester"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: make([]*ast.Field, 0),
					},
				},
			},
		},
	})

	return decls
}
