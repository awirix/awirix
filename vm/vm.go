package vm

import (
	"context"
	"github.com/spf13/viper"
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/key"
	"github.com/vivi-app/vivi/log"
	"github.com/vivi-app/vivi/lualib"
)

type Options struct {
	Silent  bool
	Context context.Context
}

func New(options *Options) *lua.LState {
	if options == nil {
		options = &Options{}
	}

	libs := []lua.LGFunction{
		lua.OpenBase,
		lua.OpenTable,
		lua.OpenString,
		lua.OpenMath,
		lua.OpenCoroutine,
		lua.OpenChannel,
		lua.OpenPackage,
	}

	luaOptions := &lua.Options{
		SkipOpenLibs: true,
		SafeMode:     viper.GetBool(key.ExtensionsSafeMode),
	}

	if !luaOptions.SafeMode {
		libs = append(libs, lua.OpenIo, lua.OpenOs)
	}

	if options.Silent {
		luaOptions.Stdout = &log.Writer{}
	}

	L := lua.NewState(luaOptions)

	for _, lib := range libs {
		lib(L)
	}

	lualib.Preload(L)

	if options.Context != nil {
		L.SetContext(options.Context)
	}

	return L
}
