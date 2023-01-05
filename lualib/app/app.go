package app

import (
	"github.com/spf13/viper"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/luadoc"
	"runtime"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "app",
		Description: "Information about platform and configuration",
		Vars: []*luadoc.Var{
			{
				Name:        "version",
				Description: "Vivi version",
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
