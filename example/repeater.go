package example

// Repeater know how to repeat tasks with a rationale
type Repeater interface {
	Repeat(task, rationale string) (count int, err error)
}
