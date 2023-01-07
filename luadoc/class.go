package luadoc

import "github.com/awirix/lua"

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
