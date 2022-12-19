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
	Silent     bool
	WorkingDir string
	Context    context.Context
}

func New(options Options) *lua.LState {
	libs := []lua.LGFunction{
		lua.OpenBase,
		lua.OpenTable,
		lua.OpenString,
		lua.OpenMath,
		lua.OpenCoroutine,
		lua.OpenChannel,
		lua.OpenPackage,
		lua.OpenIo,
	}

	luaOptions := &lua.Options{
		SkipOpenLibs: true,
		WorkingDir:   options.WorkingDir,
		IsolateIO:    viper.GetBool(key.ExtensionsSafeMode),
	}

	if !viper.GetBool(key.ExtensionsSafeMode) {
		libs = append(libs, lua.OpenOs)
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
