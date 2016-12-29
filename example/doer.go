package example

// Doer does things, sometimes graciously
type Doer interface {
	DoIt(task string, graciously bool) (int, error)
}

// Delegater employs a Doer to complete tasks
type Delegater struct {
	Delegate Doer
}

// DoSomething passes the work to Doer
func (d *Delegater) DoSomething(task string) (int, error) {
	return d.Delegate.DoIt(task, false)
}
