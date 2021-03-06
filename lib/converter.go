package fm

import (
	"fmt"
	"go/ast"
)

const (
	spyPrefix = "Spy"
)

// SpyStructConverter converts interfaces into spies, i.e., test doubles.
// Meant to be used in conjunction with SpyFuncImplementer
type SpyStructConverter struct{}

// Convert mutates the ast.TypeSpec into a struct type with public properties
// for all parameters and all return values declared in the interface
func (s *SpyStructConverter) Convert(t *ast.TypeSpec, i *ast.InterfaceType) *ast.TypeSpec {
	var list []*ast.Field
	list = append(list, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("mu")},
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("sync"),
			Sel: ast.NewIdent("Mutex"),
		},
	})

	for _, field := range i.Methods.List {
		funcType, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		if len(field.Names) != 1 {
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
		if funcType.Results != nil && len(funcType.Results.List) > 0 {
			outputStruct := buildStruct(methodName+outputSuffix, retPrefix, funcType.Results.List)
			list = append(list, outputStruct)
		}
	}

	return &ast.TypeSpec{
		Name: ast.NewIdent(spyPrefix + t.Name.Name),
		Type: &ast.StructType{
			Fields: &ast.FieldList{List: list},
		},
	}
}

// buildInputStruct writes a struct type whose fields
// reflect the various input arguments defined in the interface
func buildStruct(fieldname, prefix string, list []*ast.Field) *ast.Field {
	var (
		fields    []*ast.Field
		argOffset int
		argName   string
	)

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
