package luadoc

import "github.com/awirix/lua"

type Var struct {
	Name        string
	Description string
	Value       lua.LValue
	Type        string
}
