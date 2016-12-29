package main

import (
	genspy "github.com/enocom/genspy/lib"
)

func main() {
	g := genspy.NewGenerator(".", "spy_test.go")
	g.Generate()
}
