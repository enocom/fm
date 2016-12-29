package example

// Doer does things, sometimes graciously
type Doer interface {
	DoIt(task string, graciously bool) (int, error)
}

// Delegater employs a Doer to complete tasks
type Delegator struct {
	Delegate Doer
	Repeater
}

// DoSomething passes the work to Doer
func (d *Delegator) DoSomething(task string) (int, error) {
	return d.Delegate.DoIt(task, false)
}

// DoSomethingAgain ensures a task is repeated with a rationale
func (d *Delegator) DoSomethingAgain(task, rationale string) (int, error) {
	return d.Repeat(task, rationale)
}
