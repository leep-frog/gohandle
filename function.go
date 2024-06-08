package gohandle

type Function interface {
	Name() string
	Func() any
}

func NewFunction(name string, f any) Function {
	return &simpleFunction{name, f}
}

type simpleFunction struct {
	name string
	f    any
}

func (sf *simpleFunction) Name() string {
	return sf.name
}

func (sf *simpleFunction) Func() any {
	return sf.f
}

func Mod() Function {
	return NewFunction("mod", func(a, b int) int {
		return a % b
	})
}

func Plus() Function {
	return NewFunction("plus", func(a, b int) int {
		return a + b
	})
}
