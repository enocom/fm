package main

import (
	genspy "github.com/enocom/genspy/lib"
)

func main() {
	conv := &genspy.SpyStructConverter{}
	g := &genspy.SpyGenerator{Conv: conv}

	c := &genspy.Cmd{Wd: ".", Dst: "spy_test.go", Gen: g}
	c.Run()
}
