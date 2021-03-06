package fm

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// SrcFileParser parses the ASTs of source files only
type SrcFileParser struct{}

// ParseDir returns AST representations of all source files (excluding test files)
// within a directory.
func (s *SrcFileParser) ParseDir(dir string) (map[string]*ast.Package, error) {
	return parser.ParseDir(token.NewFileSet(), dir, isSrcFile, parser.ParseComments)
}

// isSrcFile is an ast.Filter which removes all test files
func isSrcFile(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}
