package app

import (
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
	"runtime"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"version": version,
		"os":      os,
		"arch":    arch,
		"config":  config,
	})
}

func version(L *lua.LState) int {
	L.Push(lua.LString(constant.Version))
	return 1
}

func os(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOOS))
	return 1
}

func arch(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOARCH))
	return 1
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
