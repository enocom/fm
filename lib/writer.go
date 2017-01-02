package fm

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

// FileWriter writes the ast.File to the provided filename
type FileWriter interface {
	Write(file *ast.File, filename string) error
}

// DiskFileWriter formats an ast.File as a standard go file
// and writes it to disk
type DiskFileWriter struct{}

// Write outputs the ast.File to a file on disk specified by filename
func (d *DiskFileWriter) Write(file *ast.File, filename string) error {
	spyFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	return format.Node(spyFile, token.NewFileSet(), file)
}
