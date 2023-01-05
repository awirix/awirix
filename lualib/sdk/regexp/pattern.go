package regexp

import (
	lua "github.com/vivi-app/lua"
	"regexp"
)

const regexpTypeName = "regexp"

func pushRegexp(L *lua.LState, regexp *regexp.Regexp) {
	ud := L.NewUserData()
	ud.Value = regexp
	L.SetMetatable(ud, L.GetTypeMetatable(regexpTypeName))
	L.Push(ud)
}

func checkRegexp(L *lua.LState, n int) *regexp.Regexp {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*regexp.Regexp); ok {
		return v
	}
	L.ArgError(n, "regexp expected")
	return nil
}

func regexpFind(L *lua.LState) int {
	re := checkRegexp(L, 1)
	text := L.CheckString(2)
	matches := re.FindStringSubmatch(text)
	tbl := L.NewTable()
	for _, match := range matches {
		tbl.Append(lua.LString(match))
	}
	L.Push(tbl)
	return 1
}

func regexpMatch(L *lua.LState) int {
	re := checkRegexp(L, 1)
	text := L.CheckString(2)
	matched := re.MatchString(text)
	L.Push(lua.LBool(matched))
	return 1
}

func regexpReplace(L *lua.LState) int {
	re := checkRegexp(L, 1)
	text := L.CheckString(2)
	replacement := L.CheckString(3)
	result := re.ReplaceAllString(text, replacement)
	L.Push(lua.LString(result))
	return 1
}

func regexpSplit(L *lua.LState) int {
	re := checkRegexp(L, 1)
	text := L.CheckString(2)
	tbl := L.NewTable()
	for _, match := range re.Split(text, -1) {
		tbl.Append(lua.LString(match))
	}
	L.Push(tbl)
	return 1
}

func regexpGroups(L *lua.LState) int {
	re := checkRegexp(L, 1)
	value := L.CheckString(2)

	// match all groups as map
	matches := re.FindStringSubmatch(value)
	if matches == nil {
		L.Push(lua.LNil)
		return 1
	}

	tbl := L.NewTable()
	for i, name := range re.SubexpNames() {
		if i == 0 {
			continue
		}
		tbl.RawSetString(name, lua.LString(matches[i]))
	}

	L.Push(tbl)
	return 1
}
