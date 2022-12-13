package js

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	registerVMType(L)
	registerVMValueType(L)

	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"new_vm": newVM,
	})
}
