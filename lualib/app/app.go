package app

import (
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
	"runtime"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, map[string]lua.LValue{
		"version": lua.LString(constant.Version),
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
