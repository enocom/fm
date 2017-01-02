package fm_test

import (
	"errors"
	"go/ast"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	fm "github.com/enocom/fm/lib"
)

// TestRunWritesToFile is an integration test which ensures Cmd.Run
// reads from disk, generates spies, and writes the result back out
func TestRunWritesToFile(t *testing.T) {
	wd, err, rmTempFile := writeTmpFile("package sample")
	defer rmTempFile()
	if err != nil {
		t.Fatalf("writeTmpFile failed with %v", err)
	}

	cmd := &fm.Cmd{
		DeclGenerator: buildGen(),
		Parser:        &fm.SrcFileParser{},
		FileWriter:    &fm.DiskFileWriter{},
	}
	cmd.Run(wd, "sample_test.go")

	f, err := os.Open(path.Join(wd, "sample_test.go"))
	if err != nil {
		t.Fatalf("open sample_test.go failed with %v", err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("ReadAll failed with %v", err)
	}

	want := "package sample_test"
	got := strings.TrimSpace(string(bytes))

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// TestRunReturnsErrorWhenParseFails ensures the parse error is
// returned to the caller
func TestRunReturnsErrorWhenParseFails(t *testing.T) {
	spyParser := &SpyParser{}
	expectedError := errors.New("parse failed")
	spyParser.ParseDir_Output.Ret1 = expectedError
	cmd := &fm.Cmd{
		Parser:        spyParser,
		DeclGenerator: nil,
		FileWriter:    nil,
	}

	err := cmd.Run("", "sample_test.go")

	want := expectedError
	got := err

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// TestRunReturnsErrorWhenWriteFails ensures any failure of writing
// immediately results in an error return value
func TestRunReturnsErrorWhenWriteFails(t *testing.T) {
	spyParser := &SpyParser{}
	spyParser.ParseDir_Output.Ret0 = map[string]*ast.Package{
		"bogus": &ast.Package{
			Name:  "bogus",
			Files: make(map[string]*ast.File),
		},
	}
	spyParser.ParseDir_Output.Ret1 = nil
	spyFileWriter := &SpyFileWriter{}
	expectedError := errors.New("write failed")
	spyFileWriter.Write_Output.Ret0 = expectedError

	cmd := &fm.Cmd{
		Parser:        spyParser,
		DeclGenerator: nil,
		FileWriter:    spyFileWriter,
	}

	err := cmd.Run("", "sample_test.go")

	want := expectedError
	got := err

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func writeTmpFile(code string) (string, error, func()) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err, func() {}
	}
	rmDir := func() { os.Remove(dir) }

	f, err := os.Create(path.Join(dir, "sample.go"))
	defer f.Close()
	if err != nil {
		return "", err, rmDir
	}

	_, err = f.Write([]byte(code))
	if err != nil {
		return "", err, rmDir
	}

	return dir, err, rmDir
}

func buildGen() fm.DeclGenerator {
	return &fm.SpyGenerator{
		Conv: &fm.SpyStructConverter{},
		Impl: &fm.SpyFuncImplementer{},
	}
}
