package fm

import (
	"go/ast"
	"path"
	"strings"
)

const (
	inputSuffix  = "_Input"
	outputSuffix = "_Output"
	argPrefix    = "Arg"
	retPrefix    = "Ret"
)

// DeclGenerator creates a new slice of ast declarations based on the input
type DeclGenerator interface {
	Generate(ds []ast.Decl) []ast.Decl
}

// Parser is responsible for returning the ASTs of all files
// within a directory
type Parser interface {
	ParseDir(dir string) (map[string]*ast.Package, error)
}

// Writer writes the ast.File to the provided filename
type Writer interface {
	Write(file *ast.File, filename string) error
}

// ImportWriter runs goimport against the specified file, writing the results
// back out to the same file
type ImportWriter interface {
	Write(filename string) error
}

// Cmd coordinates between a parser, a generator, and a file writer.
// It passes a parsed AST to the generator which produces an AST of spies
// from the original AST, and then passes the generated AST to the file
// writer, which saves the result to disk in the form of regular Go code.
type Cmd struct {
	DeclGenerator
	Parser
	Writer
	ImportWriter
}

// Run parses the AST within the working directory and passes it to
// the declaration generator. The result of the generator is then written
// to the designated destination with *_test as the new package name
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

		if !strings.HasSuffix(outputFilename, ".go") {
			outputFilename += ".go"
		}

		filename := path.Join(directory, outputFilename)
		err = c.Writer.Write(astFile, filename)
		if err != nil {
			return err
		}

		err = c.ImportWriter.Write(filename)
		if err != nil {
			return err
		}
	}

	return nil
}
