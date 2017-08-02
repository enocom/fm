package fm

import "os/exec"

// GoImportsWriter uses the goimports command line tool to add import statements
// to a Go source file
type GoImportsWriter struct{}

// Write passes the specified filanme through to the goimports tool
func (*GoImportsWriter) Write(filename string) error {
	cmd := exec.Command("goimports", "-w", filename)
	return cmd.Run()
}
