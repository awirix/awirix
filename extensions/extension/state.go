package extension

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/key"
	"github.com/vivi-app/vivi/log"
	"github.com/vivi-app/vivi/lualib"
	"github.com/vivi-app/vivi/luautil"
	"github.com/vivi-app/vivi/scraper"
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
	inject(e, L)

	for _, lib := range libs {
		lib(L)
	}

	lualib.Preload(L)

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

func inject(ext *Extension, L *lua.LState) {
	table := luautil.NewTable(L, map[string]lua.LValue{
		"path": lua.LString(ext.Path()),
		"passport": luautil.NewTable(L, map[string]lua.LValue{
			"name":    lua.LString(ext.Passport().Name),
			"version": lua.LString(ext.Passport().Version().String()),
			"about":   lua.LString(ext.Passport().About),
			"config": luautil.NewTable(L, nil, map[string]lua.LGFunction{
				"get": func(L *lua.LState) int {
					key := L.ToString(1)
					value, ok := ext.Passport().Config[key]
					if !ok {
						L.ArgError(1, fmt.Sprintf(`config key "%s" not found`, key))
						return 1
					}

					lvalue, err := luautil.ToLValue(L, value)
					if err != nil {
						L.RaiseError(err.Error())
					}

					L.Push(lvalue)
					return 1
				},
				"has": func(L *lua.LState) int {
					key := L.ToString(1)
					_, ok := ext.Passport().Config[key]
					L.Push(lua.LBool(ok))
					return 1
				},
			}),
		}, nil),
	}, nil)

	L.SetGlobal(scraper.GlobalExtension, table)
}
