package passport

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/extensions"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	ext := L.Context().Value(true).(extensions.ExtensionContainer)

	return luautil.NewTable(L, map[string]lua.LValue{
		"path":    lua.LString(ext.Path()),
		"version": lua.LString(ext.Passport().Version().String()),
	}, map[string]lua.LGFunction{
		"config": passportConfig,
	})
}

func getPassportFromCtx(L *lua.LState) *passport.Passport {
	return L.Context().Value(true).(extensions.ExtensionContainer).Passport()
}

func passportConfig(L *lua.LState) int {
	p := getPassportFromCtx(L)
	key := L.CheckString(1)

	section, ok := p.Config[key]
	if !ok {
		L.RaiseError("passport config section %s not found", key)
		return 1
	}

	value := section.Value()
	lvalue, err := luautil.ToLValue(L, value)
	if err != nil {
		L.RaiseError(err.Error())
		return 1
	}

	L.Push(lvalue)
	return 1
}
