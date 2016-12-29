package example_test

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
}

func (f *FakeDoer) DoIt(task string, graciously bool) (int, error) {
	f.DoIt_Called = true
	f.DoIt_Input.Arg0 = task
	f.DoIt_Input.Arg1 = graciously
	return f.DoIt_Output.Ret0, f.DoIt_Output.Ret1
}

type FakeRepeater struct {
	Repeat_Called bool
	Repeat_Input  struct {
		Arg0 string
		Arg1 string
	}
	Repeat_Output struct {
		Ret0 int
		Ret1 error
	}
}

func (f *FakeRepeater) Repeat(task, rationale string) (count int, err error) {
	f.Repeat_Called = true
	f.Repeat_Input.Arg0 = task
	f.Repeat_Input.Arg1 = rationale
	return f.Repeat_Output.Ret0, f.Repeat_Output.Ret1
}
