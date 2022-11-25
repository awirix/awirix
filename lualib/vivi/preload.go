package vivi

import (
	"github.com/vivi-app/vivi/constant"
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule(constant.App, Loader)
}

func Loader(L *lua.LState) int {
	t := L.NewTable()

	t.RawSet(lua.LString("api"), NewTableWithFuncs(L, api))
	t.RawSet(lua.LString("meta"), NewTableWithFuncs(L, meta))

	L.Push(t)
	return 1
}

func NewTableWithFuncs(L *lua.LState, funcMap map[string]lua.LGFunction) *lua.LTable {
	t := L.NewTable()
	L.SetFuncs(t, funcMap)
	return t
}

var api = map[string]lua.LGFunction{
	"watch": Watch,
}

var meta = map[string]lua.LGFunction{
	"version": Version,
	"os":      OS,
	"arch":    Arch,
}
