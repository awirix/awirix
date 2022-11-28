package regexp

import (
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
	"regexp"
)

func New(L *lua.LState) *lua.LTable {
	registerRegexpType(L)
	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"match":   match,
		"compile": compile,
	})
}

func match(L *lua.LState) int {
	pattern := L.CheckString(1)
	text := L.CheckString(2)
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(matched))
	return 1
}

func compile(L *lua.LState) int {
	pattern := L.CheckString(1)
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushRegexp(L, compiled)
	return 1
}
