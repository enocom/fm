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

func TestConvertAddsMutexWithRecordOfFunctionCallAsFields(t *testing.T) {
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

	mutexField := structType.Fields.List[0]
	want := "mu"
	got := mutexField.Names[0].Name
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	calledField := structType.Fields.List[1]
	want = "Test_Called"
	got = calledField.Names[0].Name
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// does not add input when there are no arguments
// and does not add output struct when there are no result values
func TestConvertGeneratesInputStructFieldWithArguments(t *testing.T) {
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

	want := 2
	got := len(structType.Fields.List)
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// generates Input struct with arguments
func TestConvertGeneratesInputStructWithArgumentFields(t *testing.T) {
	converter := &fm.SpyStructConverter{}

	typeSpec := converter.Convert(
		&ast.TypeSpec{Name: ast.NewIdent("Tester")},
		&ast.InterfaceType{
			Methods: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{ast.NewIdent("Test")},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{ast.NewIdent("foobar")},
										Type:  ast.NewIdent("string"),
									},
								},
							},
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

	want := 3 // mu, Test_Called, and Test_Input
	got := len(structType.Fields.List)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	inputStruct := structType.Fields.List[2]

	wantName := "Test_Input"
	gotName := inputStruct.Names[0].Name

	if wantName != gotName {
		t.Errorf("want %v, got %v", wantName, gotName)
	}

	input, ok := inputStruct.Type.(*ast.StructType)
	if !ok {
		t.Fatal("expected inputStruct to be of type StructType")
	}

	want = 1
	got = len(input.Fields.List)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	wantName = "Arg0"
	gotName = input.Fields.List[0].Names[0].Name

	if wantName != gotName {
		t.Errorf("want %v, got %v", wantName, gotName)
	}

	ident, ok := input.Fields.List[0].Type.(*ast.Ident)
	if !ok {
		t.Fatal("expected first input field to be of type *ast.Ident")
	}

	wantType := "string"
	gotType := ident.Name

	if wantType != gotType {
		t.Errorf("want %v, got %v", wantType, gotType)
	}
}

// generates Input struct with arguments when declared together
// e.g., Foobar(one, two string)
func TestConvertHandlesArgumentsDeclaredWithOneType(t *testing.T) {
	converter := &fm.SpyStructConverter{}

	typeSpec := converter.Convert(
		&ast.TypeSpec{Name: ast.NewIdent("Tester")},
		&ast.InterfaceType{
			Methods: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{ast.NewIdent("Test")},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{ast.NewIdent("foo"), ast.NewIdent("baz")},
										Type:  ast.NewIdent("string"),
									},
								},
							},
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
	inputStruct := structType.Fields.List[2]
	input, ok := inputStruct.Type.(*ast.StructType)
	if !ok {
		t.Fatal("expected inputStruct to be of type StructType")
	}

	wantName0 := "Arg0"
	gotName0 := input.Fields.List[0].Names[0].Name
	if wantName0 != gotName0 {
		t.Errorf("want %v, got %v", wantName0, gotName0)
	}

	wantName1 := "Arg1"
	gotName1 := input.Fields.List[1].Names[0].Name
	if wantName1 != gotName1 {
		t.Errorf("want %v, got %v", wantName1, gotName1)
	}

	ident0, ok := input.Fields.List[0].Type.(*ast.Ident)
	if !ok {
		t.Fatal("expected first input field to be of type *ast.Ident")
	}

	wantType0 := "string"
	gotType0 := ident0.Name

	if wantType0 != gotType0 {
		t.Errorf("want %v, got %v", wantType0, gotType0)
	}

	ident1, ok := input.Fields.List[1].Type.(*ast.Ident)
	if !ok {
		t.Fatal("expected first input field to be of type *ast.Ident")
	}

	wantType1 := "string"
	gotType1 := ident1.Name

	if wantType1 != gotType1 {
		t.Errorf("want %v, got %v", wantType1, gotType1)
	}
}

// generates Output struct with result values
func TestConvertGeneratesOutputStruct(t *testing.T) {
	converter := &fm.SpyStructConverter{}

	typeSpec := converter.Convert(
		&ast.TypeSpec{Name: ast.NewIdent("Tester")},
		&ast.InterfaceType{
			Methods: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{ast.NewIdent("Test")},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
							Results: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{ast.NewIdent("foobar")},
										Type:  ast.NewIdent("string"),
									},
								},
							},
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

	want := 3 // mu, Test_Called, and Test_Output
	got := len(structType.Fields.List)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	outputStruct := structType.Fields.List[2]

	wantName := "Test_Output"
	gotName := outputStruct.Names[0].Name

	if wantName != gotName {
		t.Errorf("want %v, got %v", wantName, gotName)
	}

	output, ok := outputStruct.Type.(*ast.StructType)
	if !ok {
		t.Fatal("expected outputStruct to be of type StructType")
	}

	want = 1
	got = len(output.Fields.List)

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	wantName = "Ret0"
	gotName = output.Fields.List[0].Names[0].Name

	if wantName != gotName {
		t.Errorf("want %v, got %v", wantName, gotName)
	}

	ident, ok := output.Fields.List[0].Type.(*ast.Ident)
	if !ok {
		t.Fatal("expected first output field to be of type *ast.Ident")
	}

	wantType := "string"
	gotType := ident.Name

	if wantType != gotType {
		t.Errorf("want %v, got %v", wantType, gotType)
	}
}

func TestConvertIgnoresEmptyReturnValues(t *testing.T) {
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
							Results: nil,
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

	want := 2 // mu and Test_Called
	got := len(structType.Fields.List)
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// TODO: what if there are multiple named returns?
