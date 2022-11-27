package sha512

import (
	"crypto/sha512"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"sum": sum,
	})
}

func sum(L *lua.LState) int {
	value := L.CheckString(1)
	s := sha512.Sum512([]byte(value))
	L.Push(lua.LString(s[:]))
	return 1
}
