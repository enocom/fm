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
func (f *FakeDoer) DoItAgain(task, prefix string) (count int, err error) {
	f.DoItAgain_Called = true
	f.DoItAgain_Input.Arg0 = task
	f.DoItAgain_Input.Arg1 = prefix
	return f.DoItAgain_Output.Ret0, f.DoItAgain_Output.Ret1
}
