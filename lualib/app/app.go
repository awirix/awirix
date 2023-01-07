package app

import (
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/luadoc"
	lua "github.com/awirix/lua"
	"github.com/spf13/viper"
	"runtime"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "app",
		Description: "Information about platform and configuration",
		Vars: []*luadoc.Var{
			{
				Name:        "version",
				Description: app.Name + " version",
				Value:       lua.LString(app.Version),
			},
			{
				Name:        "os",
				Description: "Operating system",
				Value:       lua.LString(runtime.GOOS),
			},
			{
				Name:        "arch",
				Description: "Architecture",
				Value:       lua.LString(runtime.GOARCH),
			},
		},
		Funcs: []*luadoc.Func{
			{
				Name:        "config",
				Description: "Get configuration value",
				Value:       config,
				Params: []*luadoc.Param{
					{
						Name: "key",
						Type: "string",
					},
				},
				Returns: []*luadoc.Param{
					{
						Name: "value",
						Type: luadoc.Any,
					},
				},
			},
		},
	}
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
