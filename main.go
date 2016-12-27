package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	if err := filepath.Walk(dir, processFile); err != nil {
		fatal(err)
	}
}

func isValidFile(fname string) bool {
	return strings.HasSuffix(fname, ".go") &&
		!strings.HasSuffix(fname, "_test.go")
}

var fakeTemplate string = `package {{.Pname}}_test

type Fake{{.Iname}} struct {}
`

var t = template.Must(template.New("fake_template").Parse(fakeTemplate))

type Bindings struct {
	Pname string
	Iname string
}

func processFile(path string, info os.FileInfo, err error) error {
	// temporary
	if info.Name() == "main.go" {
		return nil
	}

	if info.IsDir() {
		return nil
	}

	if !isValidFile(info.Name()) {
		return nil
	}

	// create ast from file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, info.Name(), nil, 0)
	if err != nil {
		fatal(err)
	}

	var decls []ast.Decl
	for _, d := range f.Decls {
		if g, ok := d.(*ast.GenDecl); ok {
			for _, s := range g.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					if _, ok := ts.Type.(*ast.InterfaceType); ok {
						decls = append(decls, d)
					}
				}
			}
		}
	}
	f.Decls = decls

	err = printer.Fprint(os.Stdout, fset, f)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
