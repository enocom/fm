package main

import (
	genspy "github.com/enocom/genspy/lib"
)

func main() {
	g := &genspy.SpyGenerator{}

	c := &genspy.Cmd{Wd: ".", Dst: "spy_test.go", Gen: g}
	c.Run()
}
