package example_test

import (
	"testing"

	"github.com/enocom/genspy/example"
)

func TestDelegatorCallsDoer(t *testing.T) {
	fakeDoer := &FakeDoer{}
	d := &example.Delegater{Delegate: fakeDoer}

	d.DoSomething("laundry")

	want := true
	got := fakeDoer.DoIt_Called

	if want != got {
		t.Errorf("wanted: %v, but got %v", want, got)
	}
}

func TestDelegatorCallsDoerWithArgs(t *testing.T) {
	fakeDoer := &FakeDoer{}
	d := &example.Delegater{Delegate: fakeDoer}

	d.DoSomething("laundry")

	wantArg0 := "laundry"
	gotArg0 := fakeDoer.DoIt_Input.Arg0

	if wantArg0 != gotArg0 {
		t.Errorf("wanted: %v, but got %v", wantArg0, gotArg0)
	}

	wantArg1 := false
	gotArg1 := fakeDoer.DoIt_Input.Arg1

	if wantArg1 != gotArg1 {
		t.Errorf("wanted: %v, but got %v", wantArg1, gotArg1)
	}
}

func TestDelegatorReturnsDoerResult(t *testing.T) {
	fakeDoer := &FakeDoer{}
	fakeDoer.DoIt_Output.Ret0 = 42
	fakeDoer.DoIt_Output.Ret1 = nil
	d := &example.Delegater{Delegate: fakeDoer}

	n, err := d.DoSomething("laundry")

	wantRet0 := 42
	gotRet0 := n

	if wantRet0 != gotRet0 {
		t.Errorf("wanted: %v, but got %v", wantRet0, gotRet0)
	}

	var wantRet1 error = nil
	gotRet1 := err

	if wantRet1 != gotRet1 {
		t.Errorf("wanted: %v, but got %v", wantRet1, gotRet1)
	}
}

type FakeDoer struct {
	DoIt_Called bool
	DoIt_Input  struct {
		Arg0 string
		Arg1 bool
	}
	DoIt_Output struct {
		Ret0 int
		Ret1 error
	}
	DoItAgain_Called bool
	DoItAgain_Input  struct {
		Arg0 string
		Arg1 string
	}
	DoItAgain_Output struct {
		Ret0 int
		Ret1 error
	}
}

func (f *FakeDoer) DoIt(task string, repeat bool) (int, error) {
	f.DoIt_Called = true
	f.DoIt_Input.Arg0 = task
	f.DoIt_Input.Arg1 = repeat

	return f.DoIt_Output.Ret0, f.DoIt_Output.Ret1
}

func (f *FakeDoer) DoItAgain(task, prefix string) (int, error) {
	f.DoItAgain_Called = true
	f.DoItAgain_Input.Arg0 = task
	f.DoItAgain_Input.Arg1 = prefix

	return f.DoItAgain_Output.Ret0, f.DoItAgain_Output.Ret1
}
