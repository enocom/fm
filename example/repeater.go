package example

type Repeater interface {
	Repeat(task, prefix string) (count int, err error)
}
