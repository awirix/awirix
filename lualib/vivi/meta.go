package vivi

import (
	"github.com/vivi-app/vivi/constant"
	lua "github.com/yuin/gopher-lua"
	"runtime"
)

func Version(L *lua.LState) int {
	L.Push(lua.LString(constant.Version))
	return 1
}

func OS(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOOS))
	return 1
}

func Arch(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOARCH))
	return 1
}
