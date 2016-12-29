package main

import (
	genspy "github.com/enocom/genspy/lib"
)

func main() {
	g := &genspy.Generator{Wd: ".", Dst: "spy_test.go"}
	g.GenerateSpies()
}
