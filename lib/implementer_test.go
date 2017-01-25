package fm_test

import (
	"go/ast"
	"testing"

	fm "github.com/enocom/fm/lib"
)

func TestImplementAssignsCalledField(t *testing.T) {
	someInterface := buildInterface()
	s := &fm.SpyFuncImplementer{}
	funcDecls := s.Implement(ast.NewIdent("SomeStruct"), someInterface)

	got := len(funcDecls)
	want := 1

	if want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func buildInterface() *ast.InterfaceType {
	params := []*ast.Field{}

	methodList := []*ast.Field{{
		Names: []*ast.Ident{ast.NewIdent("SomeMethod")},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{},
		},
	}}

	fieldList := &ast.FieldList{List: methodList}
	return &ast.InterfaceType{Methods: fieldList}
}

// TODO skips non FuncTypes
// TODO checks for `foo, bar string` types
// TODO assigns s.Foo_Input.Arg0 = arg, etc.
// TODO adds return values, e.g., return Foo_Output.Ret0, etc.
// TODO doesn't add return value when there is none (doesn't blow up)
