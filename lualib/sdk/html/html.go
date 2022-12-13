package html

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	registerDocumentType(L)
	registerSelectionType(L)

	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"parse": parse,
	})
}
