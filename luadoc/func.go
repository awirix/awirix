package luadoc

import (
	"github.com/awirix/lua"
	"strings"
)

type Param struct {
	Name        string
	Description string
	Type        string
	Optional    bool
}

func (p *Param) String() string {
	return p.Name
}

type Func struct {
	Name        string
	Description string
	Value       lua.LGFunction
	Generics    []string
	Params      []*Param
	Returns     []*Param
}

func (f Func) AsType() string {
	var b strings.Builder
	b.WriteString("fun(")
	for i, param := range f.Params {
		b.WriteString(param.Name)
		b.WriteString(": ")
		b.WriteString(param.Type)
		if param.Optional {
			b.WriteString("?")
		}

		if i < len(f.Params)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(")")

	if len(f.Returns) > 0 {
		b.WriteString(": ")
		b.WriteString(f.Returns[0].Type)
	}

	return b.String()
}
