package example_test

import (
	"errors"
	"testing"

	"github.com/enocom/fm/example"
)

func TestDelegatorCallsDoer(t *testing.T) {
	spyDoer := &SpyDoer{}
	d := &example.Delegator{Delegate: spyDoer}

	d.DoSomething("laundry")

	want := true
	got := spyDoer.DoIt_Called

	if want != got {
		t.Errorf("wanted: %v, but got %v", want, got)
	}
}

func TestDelegatorCallsDoerWithArgs(t *testing.T) {
	spyDoer := &SpyDoer{}
	d := &example.Delegator{Delegate: spyDoer}

	d.DoSomething("laundry")

	wantArg0 := "laundry"
	gotArg0 := spyDoer.DoIt_Input.Arg0

	if wantArg0 != gotArg0 {
		t.Errorf("wanted: %v, but got %v", wantArg0, gotArg0)
	}

	wantArg1 := false
	gotArg1 := spyDoer.DoIt_Input.Arg1

	if wantArg1 != gotArg1 {
		t.Errorf("wanted: %v, but got %v", wantArg1, gotArg1)
	}
}

func TestDelegatorReturnsDoerResult(t *testing.T) {
	spyDoer := &SpyDoer{}
	spyDoer.DoIt_Output.Ret0 = 42
	expectedErr := errors.New("some-error")
	spyDoer.DoIt_Output.Ret1 = expectedErr
	d := &example.Delegator{Delegate: spyDoer}

	n, err := d.DoSomething("laundry")

	wantRet0 := 42
	gotRet0 := n

	if wantRet0 != gotRet0 {
		t.Errorf("wanted: %v, but got %v", wantRet0, gotRet0)
	}

	wantRet1 := expectedErr
	gotRet1 := err

	if wantRet1 != gotRet1 {
		t.Errorf("wanted: %v, but got %v", wantRet1, gotRet1)
	}
}

func TestDelegatorCallsRepeater(t *testing.T) {
	r := &SpyRepeater{}
	d := &example.Delegator{Repeater: r}

	d.DoSomethingAgain("laundry", "still not done")

	want := true
	got := r.Repeat_Called

	if want != got {
		t.Errorf("wanted %v, but got %v", want, got)
	}
}

func TestDelegatorCallsRepeaterWithArgs(t *testing.T) {
	r := &SpyRepeater{}
	d := &example.Delegator{Repeater: r}

	d.DoSomethingAgain("walk the dog", "he still won't sit still")

	wantArg0 := "walk the dog"
	gotArg0 := r.Repeat_Input.Arg0

	if wantArg0 != gotArg0 {
		t.Errorf("wanted %v, but got %v", wantArg0, gotArg0)
	}

	wantArg1 := "he still won't sit still"
	gotArg1 := r.Repeat_Input.Arg1

	if wantArg1 != gotArg1 {
		t.Errorf("wanted %v, but got %v", wantArg1, gotArg1)
	}
}

func TestDelegatorReturnsRepeaterResult(t *testing.T) {
	r := &SpyRepeater{}
	r.Repeat_Output.Ret0 = 42
	expectedErr := errors.New("cat refuses")
	r.Repeat_Output.Ret1 = expectedErr
	d := &example.Delegator{Repeater: r}

	num, err := d.DoSomethingAgain("walk the cat", "it's trying to kill me")

	wantRet0 := 42
	gotRet0 := num

	if wantRet0 != gotRet0 {
		t.Errorf("wanted %v, but got %v", wantRet0, gotRet0)
	}

	wantRet1 := expectedErr
	gotRet1 := err

	if wantRet1 != gotRet1 {
		t.Errorf("wanted %v, but got %v", wantRet1, gotRet1)
	}
}
