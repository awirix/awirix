package lualib

import (
	app2 "github.com/awirix/awirix/app"
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/awirix/lualib/api"
	"github.com/awirix/awirix/lualib/app"
	"github.com/awirix/awirix/lualib/sdk"
	lua "github.com/awirix/lua"
)

func Preload(L *lua.LState) {
	lib := Lib(L)
	L.PreloadModule(lib.Name, lib.Loader())
}

func Lib(L *lua.LState) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        app2.Name,
		Description: app2.Name + " library",
		Libs: []*luadoc.Lib{
			sdk.Lib(L),
			app.Lib(),
			api.Lib(),
		},
	}
}
