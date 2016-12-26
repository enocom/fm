//go:generate genspy
package main

type Doer interface {
	DoIt() int
}

type Fooer interface {
	Foo() int
}
