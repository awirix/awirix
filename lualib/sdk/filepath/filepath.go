package filepath

import (
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
	"path/filepath"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"join":  join,
		"split": split,
		"clean": clean,
		"base":  base,
		"dir":   dir,
		"ext":   ext,
		"abs":   abs,
		"rel":   rel,
	})
}

func join(L *lua.LState) int {
	L.Push(lua.LString(filepath.Join(L.CheckString(1), L.CheckString(2))))
	return 1
}

func split(L *lua.LState) int {
	dir, file := filepath.Split(L.CheckString(1))
	L.Push(lua.LString(dir))
	L.Push(lua.LString(file))
	return 1
}

func clean(L *lua.LState) int {
	L.Push(lua.LString(filepath.Clean(L.CheckString(1))))
	return 1
}

func base(L *lua.LState) int {
	L.Push(lua.LString(filepath.Base(L.CheckString(1))))
	return 1
}

func dir(L *lua.LState) int {
	L.Push(lua.LString(filepath.Dir(L.CheckString(1))))
	return 1
}

func ext(L *lua.LState) int {
	L.Push(lua.LString(filepath.Ext(L.CheckString(1))))
	return 1
}

func abs(L *lua.LState) int {
	abs, err := filepath.Abs(L.CheckString(1))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(abs))
	return 1
}

func rel(L *lua.LState) int {
	rel, err := filepath.Rel(L.CheckString(1), L.CheckString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(rel))
	return 1
}
