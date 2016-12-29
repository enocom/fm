package main

import (
	genspy "github.com/enocom/genspy/lib"
)

func main() {
	c := &genspy.Cmd{Wd: ".", Dst: "spy_test.go"}
	c.Run()
}
