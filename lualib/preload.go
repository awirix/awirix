package lualib

import (
	"github.com/vivi-app/vivi/lualib/vivi"
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	for _, preload := range []func(*lua.LState){
		vivi.Preload,
	} {
		preload(L)
	}
}
