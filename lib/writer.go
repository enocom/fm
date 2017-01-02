package fm

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

// ASTWriter writes the ast.File to the provided filename
type ASTWriter interface {
	Write(file *ast.File, filename string) error
}

// DiskASTWriter saves an AST to disk
type DiskASTWriter struct{}

// Write outputs the ast.File to a file on disk specified by filename
func (d *DiskASTWriter) Write(file *ast.File, filename string) error {
	spyFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	return format.Node(spyFile, token.NewFileSet(), file)
}
