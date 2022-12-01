package js

import (
	"github.com/vivi-app/vivi/luautil"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	registerVMType(L)
	registerVMValueType(L)

	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"new_vm": newVM,
	})
}
