package lualib

import (
	lua "github.com/vivi-app/lua"
	app2 "github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/luadoc"
	"github.com/vivi-app/vivi/lualib/api"
	"github.com/vivi-app/vivi/lualib/app"
	"github.com/vivi-app/vivi/lualib/sdk"
)

func Preload2(L *lua.LState) {
	lib := Lib(L)
	L.PreloadModule(lib.Name, lib.Loader())
}

func Preload(L *lua.LState) {
	L.PreloadModule(app2.Name, Loader)
}

func Lib(L *lua.LState) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "vivi",
		Description: "Vivi library",
		Libs: []*luadoc.Lib{
			sdk.Lib(L),
		},
	}
}

func Loader(L *lua.LState) int {
	libs := L.NewTable()

	for name, create := range map[string]func(*lua.LState) *lua.LTable{
		"app": app.New,
		"api": api.New,
		"sdk": sdk.New,
	} {
		L.SetField(libs, name, create(L))
	}

	L.Push(libs)
	return 1
}
