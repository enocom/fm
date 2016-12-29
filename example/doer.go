package example

// Doer does things, sometimes repeatedly
type Doer interface {
	DoIt(task string, graciously bool) (int, error)
}

type Delegater struct {
	Delegate Doer
}

func (d *Delegater) DoSomething(task string) (int, error) {
	return d.Delegate.DoIt(task, false)
}
