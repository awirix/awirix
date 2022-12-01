package strings

import (
	"github.com/vivi-app/vivi/luautil"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"replace_all":    replaceAll,
		"replace":        replace,
		"split":          split,
		"to_lower":       toLower,
		"to_upper":       toUpper,
		"trim_space":     trimSpace,
		"trim_prefix":    trimPrefix,
		"trim_suffix":    trimSuffix,
		"trim":           trim,
		"trim_left":      trimLeft,
		"trim_right":     trimRight,
		"has_prefix":     hasPrefix,
		"has_suffix":     hasSuffix,
		"contains":       contains,
		"contains_any":   containsAny,
		"count":          count,
		"count_any":      countAny,
		"equal_fold":     equalFold,
		"title":          title,
		"index":          index,
		"last_index":     lastIndex,
		"repeat":         repeat,
		"index_any":      indexAny,
		"last_index_any": lastIndexAny,
	})
}

func replaceAll(L *lua.LState) int {
	s := L.CheckString(1)
	old := L.CheckString(2)
	neww := L.CheckString(3)

	L.Push(lua.LString(strings.ReplaceAll(s, old, neww)))
	return 1
}

func replace(L *lua.LState) int {
	s := L.CheckString(1)
	old := L.CheckString(2)
	neww := L.CheckString(3)
	n := L.CheckInt(4)

	L.Push(lua.LString(strings.Replace(s, old, neww, n)))
	return 1
}

func split(L *lua.LState) int {
	s := L.CheckString(1)
	sep := L.CheckString(2)

	table, _ := luautil.ToLValue(L, strings.Split(s, sep))
	L.Push(table)
	return 1
}

func toLower(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.ToLower(s)))
	return 1
}

func toUpper(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.ToUpper(s)))
	return 1
}

func trimSpace(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.TrimSpace(s)))
	return 1
}

func trimPrefix(L *lua.LState) int {
	s := L.CheckString(1)
	prefix := L.CheckString(2)

	L.Push(lua.LString(strings.TrimPrefix(s, prefix)))
	return 1
}

func trimSuffix(L *lua.LState) int {
	s := L.CheckString(1)
	suffix := L.CheckString(2)

	L.Push(lua.LString(strings.TrimSuffix(s, suffix)))
	return 1
}

func trim(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LString(strings.Trim(s, cutset)))
	return 1
}

func trimLeft(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LString(strings.TrimLeft(s, cutset)))
	return 1
}

func trimRight(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LString(strings.TrimRight(s, cutset)))
	return 1
}

func hasPrefix(L *lua.LState) int {
	s := L.CheckString(1)
	prefix := L.CheckString(2)

	L.Push(lua.LBool(strings.HasPrefix(s, prefix)))
	return 1
}

func hasSuffix(L *lua.LState) int {
	s := L.CheckString(1)
	suffix := L.CheckString(2)

	L.Push(lua.LBool(strings.HasSuffix(s, suffix)))
	return 1
}

func contains(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LBool(strings.Contains(s, sub)))
	return 1
}

func containsAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LBool(strings.ContainsAny(s, cutset)))
	return 1
}

func count(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LNumber(strings.Count(s, sub)))
	return 1
}

func countAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LNumber(strings.Count(s, cutset)))
	return 1
}

func equalFold(L *lua.LState) int {
	s := L.CheckString(1)
	t := L.CheckString(2)

	L.Push(lua.LBool(strings.EqualFold(s, t)))
	return 1
}

func index(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LNumber(strings.Index(s, sub)))
	return 1
}

func indexAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LNumber(strings.IndexAny(s, cutset)))
	return 1
}

func lastIndex(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LNumber(strings.LastIndex(s, sub)))
	return 1
}

func lastIndexAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LNumber(strings.LastIndexAny(s, cutset)))
	return 1
}

func repeat(L *lua.LState) int {
	s := L.CheckString(1)
	count := L.CheckInt(2)

	L.Push(lua.LString(strings.Repeat(s, count)))
	return 1
}

func title(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.ToTitle(s)))
	return 1
}
