package app

import (
	"github.com/spf13/viper"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/luautil"
	"runtime"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, map[string]lua.LValue{
		"version": lua.LString(app.Version),
		"os":      lua.LString(runtime.GOOS),
		"arch":    lua.LString(runtime.GOARCH),
	}, map[string]lua.LGFunction{
		"config": config,
	})
}

func config(L *lua.LState) int {
	key := L.CheckString(1)

	switch value := viper.Get(key).(type) {
	case string:
		L.Push(lua.LString(value))
	case int:
		L.Push(lua.LNumber(value))
	case bool:
		L.Push(lua.LBool(value))
	default:
		L.Push(lua.LNil)
	}

	return 1
}
