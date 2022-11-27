package sdk

import (
	"github.com/vivi-app/vivi/lualib/sdk/crypto"
	"github.com/vivi-app/vivi/lualib/sdk/html"
	"github.com/vivi-app/vivi/lualib/sdk/json"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, map[string]lua.LValue{
		"json":   json.New(L),
		"html":   html.New(L),
		"crypto": crypto.New(L),
	}, nil)
}
