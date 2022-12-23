package extension

import (
	"github.com/spf13/viper"
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/key"
	"github.com/vivi-app/vivi/log"
	"github.com/vivi-app/vivi/lualib"
	"github.com/vivi-app/vivi/where"
	"path/filepath"
	"strings"
)

func (e *Extension) initState(debug bool) {
	if e.state != nil {
		return
	}

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
		WorkingDir:   e.Path(),
		IsolateIO:    viper.GetBool(key.ExtensionsSafeMode),
		TempDir:      where.Temp(),
	}

	if !viper.GetBool(key.ExtensionsSafeMode) {
		libs = append(libs, lua.OpenOs)
	}

	if !debug {
		luaOptions.Stdout = &log.Writer{}
	}

	L := lua.NewState(luaOptions)

	for _, lib := range libs {
		lib(L)
	}

	lualib.Preload(L)

	// this is hideous, but works
	e.context.Set("extension", e)
	L.SetContext(e.context)

	// add local files to the path
	pkg := L.GetGlobal("package").(*lua.LTable)
	paths := strings.Split(pkg.RawGetString("path").String(), ";")
	viviPaths := []string{
		filepath.Join(e.Path(), "?.lua"),
	}
	paths = append(viviPaths, paths...)
	pkg.RawSetString("path", lua.LString(strings.Join(paths, ";")))

	e.state = L
}
