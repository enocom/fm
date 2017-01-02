package fm

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

type ASTWriter interface {
	Write(file *ast.File, filename string) error
}

type DiskASTWriter struct{}

func (d *DiskASTWriter) Write(file *ast.File, filename string) error {
	spyFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	return format.Node(spyFile, token.NewFileSet(), file)
}
