package fm

import "os/exec"

type GoImportsWriter struct{}

func (*GoImportsWriter) Write(filename string) error {
	cmd := exec.Command("goimports", "-w", filename)
	return cmd.Run()
}
