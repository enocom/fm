package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
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

	fmt.Println("Processing -->", path)

	// create ast from file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, info.Name(), nil, 0)
	if err != nil {
		fatal(err)
	}
	// find all interfaces in file
	parseDecls(f.Decls)
	// generate spy implementations
	fi, err := os.Create("spy_test.go")
	if err != nil {
		fatal(err)
	}

	// render bindings to template
	data := Bindings{Pname: f.Name.Name, Iname: "Foobar"}
	if err = t.Execute(fi, data); err != nil {
		fatal(err)
	}

	return nil
}

func parseDecls(decls []ast.Decl) {
	for _, decl := range decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			parseInterfaces(d)
		default:
			// do nothing
		}
	}

}

func parseInterfaces(d *ast.GenDecl) {
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			fmt.Println("found a spec", s.Name)
			parseType(s.Type)
		default:
			// do nothing
		}
	}
}

func parseType(expr ast.Expr) {
	switch t := expr.(type) {
	case *ast.InterfaceType:
		fmt.Println("found an interface", t)
	default:
		// do nothing
	}
}

func fatal(err error) {
	fmt.Printf("Error %v\n", err)
	os.Exit(1)
}
