package genspy

import "go/ast"

// DeclGenerator creates a new slice of ast.Decl based on the input
type DeclGenerator interface {
	Generate(ds []ast.Decl) []ast.Decl
}

// SpyGenerator creates spy implementations of interface declarations
type SpyGenerator struct {
	Conv StructConverter
	Impl FuncImplementer
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

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			// TODO: stop mutating typeSpec
			g.Conv.Convert(typeSpec, interfaceType)
			// TODO: extract
			funcDecls := g.Impl.Implement(typeSpec, interfaceType)

			decls = append(decls, genDecl)
			for _, fd := range funcDecls {
				decls = append(decls, fd)
			}
		}
	}
	return decls
}