//go:generate genspy
package main

type Doer interface {
	DoIt(task string, repeat bool) (int, error)
	DoItAgain(task, prefix string) (count int, err error)
}
