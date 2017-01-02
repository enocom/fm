package fm

import (
	"fmt"
	"go/ast"
	"os"
	"path"
)

const (
	inputSuffix  = "_Input"
	outputSuffix = "_Output"
	argPrefix    = "Arg"
	retPrefix    = "Ret"
)

// Cmd coordinates between a parser and generator. It passes a parsed AST
// to the generator and then writes the generated code to disk.
type Cmd struct {
	Gen DeclGenerator
	Psr Parser
	Wrt FileWriter
}

// Run parses the ast within the working directory and passes it to
// the declaration generator. The result of the generator is then written
// to the designated destination with *_test.go as the new package name
func (c *Cmd) Run(directory, outputFilename string) {
	pkgs, err := c.Psr.ParseDir(directory)
	if err != nil {
		fatal(err) // TODO: return error
	}

	for pname, p := range pkgs {
		var decls []ast.Decl
		for _, f := range p.Files {
			spyDecls := c.Gen.Generate(f.Decls)
			decls = append(decls, spyDecls...)
		}

		astFile := &ast.File{
			Name:  ast.NewIdent(pname + "_test"),
			Decls: decls,
		}

		// TODO: ensure go extension is added only when necessary
		// TODO: return error
		c.Wrt.Write(astFile, path.Join(directory, outputFilename+".go"))
	}
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
