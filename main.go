package main

import (
	"flag"
	"fmt"

	fm "github.com/enocom/fm/lib"
)

const Version = "1.1.0"

func main() {
	printVersion := flag.Bool("version", false, "Print version and exit")
	outputFilename := flag.String(
		"out",
		"spy_test",
		"Name of output file with generated spies",
	)
	workingDir := flag.String(
		"dir",
		".",
		"Directory to search for interfaces",
	)
	flag.Parse()

	if *printVersion {
		fmt.Printf("fm version %s\n", Version)
		return
	}

	conv := &fm.SpyStructConverter{}
	impl := &fm.SpyFuncImplementer{}
	g := &fm.SpyGenerator{Conv: conv, Impl: impl}

	c := &fm.Cmd{Wd: *workingDir, Dst: *outputFilename + ".go", Gen: g}
	c.Run()
}
