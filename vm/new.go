package vm

import (
	"github.com/vivi-app/vivi/lualib"
	lua "github.com/yuin/gopher-lua"
	"path/filepath"
	"strings"
)

func New(path string) *lua.LState {
	L := lua.NewState()

	// Load the standard libraries except for the debug, io and os
	for _, openLib := range []lua.LGFunction{
		lua.OpenBase,
		lua.OpenTable,
		lua.OpenString,
		lua.OpenMath,
		lua.OpenCoroutine,
		lua.OpenChannel,
		lua.OpenPackage,
	} {
		openLib(L)
	}

	pkg := L.GetGlobal("package").(*lua.LTable)
	paths := strings.Split(pkg.RawGetString("path").String(), ";")

	viviPaths := []string{
		filepath.Join(path, "?.lua"),
	}

	paths = append(viviPaths, paths...)

	pkg.RawSetString("path", lua.LString(strings.Join(paths, ";")))

	lualib.Preload(L)
	return L
}
