package sdk

import (
	"github.com/vivi-app/vivi/lualib/sdk/cmd"
	"github.com/vivi-app/vivi/lualib/sdk/crypto"
	"github.com/vivi-app/vivi/lualib/sdk/html"
	"github.com/vivi-app/vivi/lualib/sdk/http"
	"github.com/vivi-app/vivi/lualib/sdk/js"
	"github.com/vivi-app/vivi/lualib/sdk/json"
	"github.com/vivi-app/vivi/lualib/sdk/regexp"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, map[string]lua.LValue{
		"json":   json.New(L),
		"html":   html.New(L),
		"crypto": crypto.New(L),
		"http":   http.New(L),
		"regexp": regexp.New(L),
		"js":     js.New(L),
		"cmd":    cmd.New(L),
	}, nil)
}
