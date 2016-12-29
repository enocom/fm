package main

import (
	fm "github.com/enocom/fm/lib"
)

func main() {
	conv := &fm.SpyStructConverter{}
	impl := &fm.SpyFuncImplementer{}
	g := &fm.SpyGenerator{Conv: conv, Impl: impl}

	c := &fm.Cmd{Wd: ".", Dst: "spy_test.go", Gen: g}
	c.Run()
}
