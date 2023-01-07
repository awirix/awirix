package extension

import (
	"context"
	"github.com/awirix/awirix/key"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/lualib"
	"github.com/awirix/awirix/luautil"
	"github.com/awirix/awirix/scraper"
	"github.com/awirix/awirix/where"
	"github.com/awirix/lua"
	"github.com/spf13/viper"
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
	inject(e, L)

	for _, lib := range libs {
		lib(L)
	}

	lualib.Preload(L)

	// add local files to the path
	pkg := L.GetGlobal("package").(*lua.LTable)
	paths := strings.Split(pkg.RawGetString("path").String(), ";")
	appPaths := []string{
		filepath.Join(e.Path(), "?.lua"),
	}
	paths = append(appPaths, paths...)
	pkg.RawSetString("path", lua.LString(strings.Join(paths, ";")))

	e.state = L
}

func inject(ext *Extension, L *lua.LState) {
	table := luautil.NewTable(L, map[string]lua.LValue{
		"path": lua.LString(ext.Path()),
		"passport": luautil.NewTable(L, map[string]lua.LValue{
			"name":    lua.LString(ext.Passport().Name),
			"version": lua.LString(ext.Passport().Version().String()),
			"about":   lua.LString(ext.Passport().About),
		}, nil),
	}, nil)

	L.SetGlobal(scraper.GlobalExtension, table)
	L.SetContext(context.WithValue(context.Background(), "extension", ext))
}
