package luadoc

import "github.com/vivi-app/lua"

type Method struct {
	Name        string
	Description string
	Value       lua.LGFunction
	Params      []*Param
	Returns     []*Param
}

type Class struct {
	Name        string
	Description string
	Methods     []*Method
}
