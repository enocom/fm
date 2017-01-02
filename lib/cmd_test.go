package fm_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	fm "github.com/enocom/fm/lib"
)

var sampleCode string = `package sample`

// TestRunWritesToFile is an integration test which ensures Cmd.Run
// reads from disk, generates spies, and writes the result back out
func TestRunWritesToFile(t *testing.T) {
	wd, err, rmTempFile := writeTmpFile(sampleCode)
	defer rmTempFile()
	if err != nil {
		t.Fatalf("writeTmpFile failed with %v", err)
	}

	cmd := &fm.Cmd{
		Gen: buildGen(),
		Psr: &fm.SrcFileParser{},
		Wrt: &fm.DiskFileWriter{},
	}
	cmd.Run(wd, "sample_test")

	f, err := os.Open(path.Join(wd, "sample_test.go"))
	if err != nil {
		t.Fatalf("open sample_test.go failed with %v", err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("ReadAll failed with %v", err)
	}

	want := `package sample_test`
	got := strings.TrimSpace(string(bytes))

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
