package fm

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

// FileWriter formats an ast.File as a standard go file
// and writes it to disk
type FileWriter struct{}

// Write outputs the ast.File to a file on disk specified by filename
func (w *FileWriter) Write(file *ast.File, filename string) error {
	spyFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	return format.Node(spyFile, token.NewFileSet(), file)
}
