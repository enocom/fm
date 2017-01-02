package fm

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

const (
	inputSuffix  = "_Input"
	outputSuffix = "_Output"
	argPrefix    = "Arg"
	retPrefix    = "Ret"
)

// Cmd passes all declarations found within the working directory to
// the declaration generator, and writes the output to the filename
// specified by Dst
type Cmd struct {
	Wd  string
	Dst string
	Gen DeclGenerator
}

// Run parses the ast within the working directory and passes it to
// the declaration generator. The result of the generator is then written
// to the designated destination with *_test.go as the new package name
func (c *Cmd) Run() {
	pkgs, err := parser.ParseDir(token.NewFileSet(), c.Wd, isSrcFile, 0)
	if err != nil {
		fatal(err)
	}

	for pname, p := range pkgs {
		spyFile, err := os.Create(path.Join(c.Wd, c.Dst))
		if err != nil {
			fatal(err)
		}

		var decls []ast.Decl
		for _, f := range p.Files {
			spyDecls := c.Gen.Generate(f.Decls)
			decls = append(decls, spyDecls...)
		}

		astFile := &ast.File{
			Name:  ast.NewIdent(pname + "_test"),
			Decls: decls,
		}

		format.Node(spyFile, token.NewFileSet(), astFile)
	}
}

// isSrcFile is an ast.Filter which removes all test files
func isSrcFile(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
