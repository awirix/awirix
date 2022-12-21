package api

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"watch":     watch,
		"open":      openDefault,
		"open_data": openData,
		"download":  download,
		"save":      save,
	})
}
