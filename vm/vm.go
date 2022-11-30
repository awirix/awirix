package vm

import (
	"github.com/vivi-app/vivi/lualib"
	lua "github.com/yuin/gopher-lua"
)

func New() *lua.LState {
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

	lualib.Preload(L)

	return L
}
