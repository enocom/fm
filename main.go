package main

import (
	"flag"
	"fmt"

	fm "github.com/enocom/fm/lib"
)

const Version = "1.0.0"

func main() {
	printVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *printVersion {
		fmt.Printf("fm version %s\n", Version)
		return
	}

	conv := &fm.SpyStructConverter{}
	impl := &fm.SpyFuncImplementer{}
	g := &fm.SpyGenerator{Conv: conv, Impl: impl}

	c := &fm.Cmd{Wd: ".", Dst: "spy_test.go", Gen: g}
	c.Run()
}
