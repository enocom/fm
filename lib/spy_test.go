package fm_test

import "go/ast" // manually adding this for now

type SpyStructConverter struct {
	Convert_Called bool
	Convert_Input  struct {
		Arg0 *ast.TypeSpec
		Arg1 *ast.InterfaceType
	}
	Convert_Output struct {
		Ret0 *ast.TypeSpec
	}
}

func (f *SpyStructConverter) Convert(t *ast.TypeSpec, i *ast.InterfaceType) *ast.TypeSpec {
	f.Convert_Called = true
	f.Convert_Input.Arg0 = t
	f.Convert_Input.Arg1 = i
	return f.Convert_Output.Ret0
}

type SpyDeclGenerator struct {
	Generate_Called bool
	Generate_Input  struct {
		Arg0 []ast.Decl
	}
	Generate_Output struct {
		Ret0 []ast.Decl
	}
}

func (f *SpyDeclGenerator) Generate(ds []ast.Decl) []ast.Decl {
	f.Generate_Called = true
	f.Generate_Input.Arg0 = ds
	return f.Generate_Output.Ret0
}

type SpyFuncImplementer struct {
	Implement_Called bool
	Implement_Input  struct {
		Arg0 *ast.TypeSpec
		Arg1 *ast.InterfaceType
	}
	Implement_Output struct {
		Ret0 []*ast.FuncDecl
	}
}

func (f *SpyFuncImplementer) Implement(t *ast.TypeSpec, i *ast.InterfaceType) []*ast.FuncDecl {
	f.Implement_Called = true
	f.Implement_Input.Arg0 = t
	f.Implement_Input.Arg1 = i
	return f.Implement_Output.Ret0
}

type SpyParser struct {
	ParseDir_Called bool
	ParseDir_Input  struct {
		Arg0 string
	}
	ParseDir_Output struct {
		Ret0 map[string]*ast.Package
		Ret1 error
	}
}

func (f *SpyParser) ParseDir(dir string) (map[string]*ast.Package, error) {
	f.ParseDir_Called = true
	f.ParseDir_Input.Arg0 = dir
	return f.ParseDir_Output.Ret0, f.ParseDir_Output.Ret1
}

type SpyFileWriter struct {
	Write_Called bool
	Write_Input  struct {
		Arg0 *ast.File
		Arg1 string
	}
	Write_Output struct {
		Ret0 error
	}
}

func (f *SpyFileWriter) Write(file *ast.File, filename string) error {
	f.Write_Called = true
	f.Write_Input.Arg0 = file
	f.Write_Input.Arg1 = filename
	return f.Write_Output.Ret0
}
