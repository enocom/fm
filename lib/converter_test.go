package fm_test

import (
	"go/ast"
	"testing"

	"github.com/enocom/fm/lib"
)

func TestConvertBuildsAddsSpyToTypeSpecName(t *testing.T) {
	converter := &fm.SpyStructConverter{}

	typeSpec := converter.Convert(
		&ast.TypeSpec{
			Name: ast.NewIdent("Tester"),
		},
		&ast.InterfaceType{
			Methods: &ast.FieldList{List: make([]*ast.Field, 0)},
		},
	)

	want := "SpyTester"
	got := typeSpec.Name.Name

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
