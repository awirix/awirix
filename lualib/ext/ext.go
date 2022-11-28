package passport

import (
	"github.com/vivi-app/vivi/lualib/ext/passport"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, map[string]lua.LValue{
		"passport": passport.New(L),
	}, nil)
}
