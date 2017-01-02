package fm

import (
	"go/ast"
	"path"
)

const (
	inputSuffix  = "_Input"
	outputSuffix = "_Output"
	argPrefix    = "Arg"
	retPrefix    = "Ret"
)

// Cmd coordinates between a parser, a generator, and a file writer.
// It passes a parsed AST to the generator which produces an AST of spies
// from the original AST, and then passes the generated AST to a the file
// writer, which saves the result to disk in the form of regular Go code.
type Cmd struct {
	DeclGenerator
	Parser
	FileWriter
}

// Run parses the AST within the working directory and passes it to
// the declaration generator. The result of the generator is then written
// to the designated destination with *_test.go as the new package name
func (c *Cmd) Run(directory, outputFilename string) error {
	pkgs, err := c.ParseDir(directory)
	if err != nil {
		return err
	}

	for pname, p := range pkgs {
		var decls []ast.Decl
		for _, f := range p.Files {
			spyDecls := c.Generate(f.Decls)
			decls = append(decls, spyDecls...)
		}

		astFile := &ast.File{
			Name:  ast.NewIdent(pname + "_test"),
			Decls: decls,
		}

		// TODO: ensure go extension is added only when necessary
		err = c.Write(astFile, path.Join(directory, outputFilename+".go"))
		if err != nil {
			return err
		}
	}

	return nil
}
