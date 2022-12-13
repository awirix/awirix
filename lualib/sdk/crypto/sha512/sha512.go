package sha512

import (
	"crypto/sha512"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"sum": sum,
	})
}

func sum(L *lua.LState) int {
	value := L.CheckString(1)
	s := sha512.Sum512([]byte(value))
	L.Push(lua.LString(s[:]))
	return 1
}
