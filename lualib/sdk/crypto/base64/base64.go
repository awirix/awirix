package base64

import (
	"encoding/base64"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"encode": encode,
		"decode": decode,
	})
}

func encode(L *lua.LState) int {
	value := L.CheckString(1)
	L.Push(lua.LString(base64.StdEncoding.EncodeToString([]byte(value))))
	return 1
}

func decode(L *lua.LState) int {
	value := L.CheckString(1)
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(decoded))
	return 1
}
