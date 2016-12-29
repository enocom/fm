package main

import (
	genspy "github.com/enocom/genspy/lib"
)

func main() {
	conv := &genspy.SpyStructConverter{}
	impl := &genspy.SpyFuncImplementer{}
	g := &genspy.SpyGenerator{Conv: conv, Impl: impl}

	c := &genspy.Cmd{Wd: ".", Dst: "spy_test.go", Gen: g}
	c.Run()
}
