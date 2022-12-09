package regexp

import (
	"github.com/mvdan/xurls"
	"github.com/vivi-app/vivi/luautil"
	lua "github.com/yuin/gopher-lua"
	"regexp"
)

func New(L *lua.LState) *lua.LTable {
	registerRegexpType(L)

	toLValue := func(r *regexp.Regexp) *lua.LUserData {
		ud := L.NewUserData()
		ud.Value = r
		L.SetMetatable(ud, L.GetTypeMetatable(regexpTypeName))
		return ud
	}

	return luautil.NewTable(L, map[string]lua.LValue{
		"urls_relaxed": toLValue(xurls.Relaxed),
		"urls_strict":  toLValue(xurls.Strict),
	}, map[string]lua.LGFunction{
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
