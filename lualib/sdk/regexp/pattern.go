package regexp

import (
	lua "github.com/awirix/lua"
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

func regexpFindSubmatch(L *lua.LState) int {
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

func regexpReplaceAll(L *lua.LState) int {
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

func regexpReplaceAllFunc(L *lua.LState) int {
	re := checkRegexp(L, 1)
	text := L.CheckString(2)
	replacer := L.CheckFunction(3)

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		L.Push(replacer)
		L.Push(lua.LString(match))
		L.Call(1, 1)
		return L.CheckString(-1)
	})

	L.Push(lua.LString(result))
	return 1
}
