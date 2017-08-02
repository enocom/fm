package main

import (
	"flag"
	"fmt"
	"os"

	fm "github.com/enocom/fm/lib"
)

// Version designates the currently released version of fm
const Version = "1.1.0"

func main() {
	printVersion := flag.Bool("version", false, "Print version and exit")
	outputFilename := flag.String(
		"out",
		"fm_test.go",
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

	c := &fm.Cmd{
		DeclGenerator: &fm.SpyGenerator{
			Converter:   &fm.SpyStructConverter{},
			Implementer: &fm.SpyFuncImplementer{},
		},
		Parser:       &fm.SrcFileParser{},
		Writer:       &fm.FileWriter{},
		ImportWriter: &fm.GoImportsWriter{},
	}

	err := c.Run(*workingDir, *outputFilename)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}
}
