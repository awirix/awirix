package lualib

import (
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/lualib/api"
	"github.com/vivi-app/vivi/lualib/app"
	ext "github.com/vivi-app/vivi/lualib/ext"
	"github.com/vivi-app/vivi/lualib/sdk"
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule(constant.App, Loader)
}

func Loader(L *lua.LState) int {
	libs := L.NewTable()

	for name, create := range map[string]func(*lua.LState) *lua.LTable{
		"app": app.New,
		"api": api.New,
		"sdk": sdk.New,
		"ext": ext.New,
	} {
		L.SetField(libs, name, create(L))
	}

	L.Push(libs)
	return 1
}
