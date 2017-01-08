package fm_test

import (
	"go/ast"
	"go/token"
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

// e.g., ast.FuncDecl
func TestGenerateSkipsDeclsThatAreNotGenDecls(t *testing.T) {
	gen := &fm.SpyGenerator{
		Converter:   &fm.SpyStructConverter{},
		Implementer: &fm.SpyFuncImplementer{},
	}
	decls := buildFuncDeclAST()
	spyDecls := gen.Generate(decls)

	want := 0
	got := len(spyDecls)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// e.g., ast.ValueSpec
func TestGenerateSkipsSpecsThatAreNotTypeSpecs(t *testing.T) {
	gen := &fm.SpyGenerator{
		Converter:   &fm.SpyStructConverter{},
		Implementer: &fm.SpyFuncImplementer{},
	}
	decls := buildValueSpecAST()
	spyDecls := gen.Generate(decls)

	want := 0
	got := len(spyDecls)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// e.g., ast.StructType
func TestGenerateSkipsTypeSpecsThatAreNotInterfaceTypes(t *testing.T) {
	gen := &fm.SpyGenerator{
		Converter:   &fm.SpyStructConverter{},
		Implementer: &fm.SpyFuncImplementer{},
	}
	decls := buildStructAST()
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

// buildTestAST generates as AST of the following code:
//
// type Tester interface {
//     Test()
// }
func buildInterfaceAST() []ast.Decl {
	var fields []*ast.Field
	fields = append(fields, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("Test")},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: nil},
			Results: &ast.FieldList{List: nil},
		},
	})

	var decls []ast.Decl
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

// buildStructAST generates as AST of the following code:
//
// type Tester struct {}
func buildStructAST() []ast.Decl {
	var decls []ast.Decl
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

// buildStructAST generates an AST of the followign code:
//
// func foobar() {}
func buildFuncDeclAST() []ast.Decl {
	var decls []ast.Decl
	decls = append(decls, &ast.FuncDecl{
		Recv: nil,
		Name: ast.NewIdent("foobar"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: make([]*ast.Field, 0),
			},
			Results: nil,
		},
		Body: &ast.BlockStmt{
			List: nil,
		},
	})

	return decls
}

// buildValueSpecAST generates an AST of the following code:
//
// var foobar string
func buildValueSpecAST() []ast.Decl {
	var decls []ast.Decl
	decls = append(decls, &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent("foobar"),
					ast.NewIdent("string"),
				},
			},
		},
	})

	return decls
}
