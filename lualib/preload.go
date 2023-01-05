package lualib

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luadoc"
	"github.com/vivi-app/vivi/lualib/api"
	"github.com/vivi-app/vivi/lualib/app"
	"github.com/vivi-app/vivi/lualib/sdk"
)

func Preload(L *lua.LState) {
	lib := Lib(L)
	L.PreloadModule(lib.Name, lib.Loader())
}

func Lib(L *lua.LState) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "vivi",
		Description: "Vivi library",
		Libs: []*luadoc.Lib{
			sdk.Lib(L),
			app.Lib(),
			api.Lib(),
		},
	}
}
