package fm

import "go/ast"

// DeclGenerator creates a new slice of ast declarations based on the input
type DeclGenerator interface {
	Generate(ds []ast.Decl) []ast.Decl
}

// SpyGenerator creates spy implementations of interface declarations
type SpyGenerator struct {
	Converter   StructConverter
	Implementer FuncImplementer
}

// Generate transforms all the interfaces in the list of declarations
// into spies in the form of structs with implemented functions
func (g *SpyGenerator) Generate(ds []ast.Decl) []ast.Decl {
	var decls []ast.Decl
	for _, d := range ds {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if len(genDecl.Specs) != 1 {
			continue // TODO: would this ever happen?
		}

		spec := genDecl.Specs[0]
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		structTypeSpec := g.Converter.Convert(typeSpec, interfaceType)
		decls = append(decls, &ast.GenDecl{
			Tok:   genDecl.Tok,
			Specs: []ast.Spec{structTypeSpec},
		})

		funcDecls := g.Implementer.Implement(structTypeSpec, interfaceType)
		for _, fd := range funcDecls {
			decls = append(decls, fd)
		}
	}

	return decls
}
