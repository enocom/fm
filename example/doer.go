//go:generate genspy
package example

// Doer does things, sometimes repeatedly
type Doer interface {
	DoIt(task string, repeat bool) (int, error)
	DoItAgain(task, prefix string) (count int, err error)
}

type Delegater struct {
	Delegate Doer
}

func (d *Delegater) DoSomething(task string) (int, error) {
	return d.Delegate.DoIt(task, false)
}

// RealDoer get stuff done for real
type RealDoer struct {
}

// DoIt does all the hard work
func (r *RealDoer) DoIt(task string, repeat bool) (int, error) {
	return 0, nil
}

const someConst = 1
