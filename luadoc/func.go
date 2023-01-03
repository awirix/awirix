package luadoc

import (
	"github.com/vivi-app/lua"
)

type Param struct {
	Name        string
	Description string
	Type        string
	Opt         bool
}

func (p *Param) String() string {
	return p.Name
}

type Func struct {
	Name        string
	Description string
	Value       lua.LGFunction
	Params      []*Param
	Returns     []*Param
}
