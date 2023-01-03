package luadoc

import "github.com/vivi-app/lua"

type Var struct {
	Name        string
	Description string
	Value       lua.LValue
}

func (v *Var) LuaDoc() string {
	return ""
}
