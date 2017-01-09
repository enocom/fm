package fm_test

import (
	"go/ast"
	"testing"

	"github.com/enocom/fm/lib"
)

func TestConvertBuildsAddsSpyToTypeSpecName(t *testing.T) {
	converter := &fm.SpyStructConverter{}

	typeSpec := converter.Convert(
		&ast.TypeSpec{Name: ast.NewIdent("Tester")},
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

func TestConvertAddsRecordOfFunctionCallAsField(t *testing.T) {
	converter := &fm.SpyStructConverter{}

	typeSpec := converter.Convert(
		&ast.TypeSpec{Name: ast.NewIdent("Tester")},
		&ast.InterfaceType{
			Methods: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{ast.NewIdent("Test")},
						Type: &ast.FuncType{
							Params:  &ast.FieldList{},
							Results: &ast.FieldList{},
						},
					},
				},
			},
		},
	)

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		t.Fatal("expected typeSpec to be of type StructType")
	}

	want := 1
	got := len(structType.Fields.List)
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	calledField := structType.Fields.List[0]

	wantName := "Test_Called"
	gotName := calledField.Names[0].Name
	if wantName != gotName {
		t.Errorf("want %v, got %v", wantName, gotName)
	}
}
