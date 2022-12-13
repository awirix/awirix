package lualib

import (
	lua "github.com/vivi-app/lua"
	app2 "github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/lualib/api"
	"github.com/vivi-app/vivi/lualib/app"
	ext "github.com/vivi-app/vivi/lualib/ext"
	"github.com/vivi-app/vivi/lualib/sdk"
)

func Preload(L *lua.LState) {
	L.PreloadModule(app2.Name, Loader)
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
