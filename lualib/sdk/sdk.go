package sdk

import (
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/awirix/lualib/sdk/cmd"
	"github.com/awirix/awirix/lualib/sdk/crypto"
	"github.com/awirix/awirix/lualib/sdk/filepath"
	"github.com/awirix/awirix/lualib/sdk/fmt"
	"github.com/awirix/awirix/lualib/sdk/functional"
	"github.com/awirix/awirix/lualib/sdk/fuzzy"
	"github.com/awirix/awirix/lualib/sdk/html"
	"github.com/awirix/awirix/lualib/sdk/http"
	"github.com/awirix/awirix/lualib/sdk/js"
	"github.com/awirix/awirix/lualib/sdk/json"
	"github.com/awirix/awirix/lualib/sdk/levenshtein"
	"github.com/awirix/awirix/lualib/sdk/pdf"
	"github.com/awirix/awirix/lualib/sdk/regexp"
	"github.com/awirix/awirix/lualib/sdk/strings"
	"github.com/awirix/awirix/lualib/sdk/time"
	"github.com/awirix/awirix/lualib/sdk/urls"
	"github.com/awirix/awirix/luautil"
	lua "github.com/awirix/lua"
)

func Lib(L *lua.LState) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "sdk",
		Description: app.Name + ` SDK library. Contains various utilities for making HTTP requests, working with JSON, HTML, and more.`,
		Libs: []*luadoc.Lib{
			regexp.Lib(L),
			strings.Lib(),
			pdf.Lib(),
			json.Lib(),
			cmd.Lib(),
			crypto.Lib(L),
			filepath.Lib(),
			js.Lib(),
			http.Lib(),
			html.Lib(),
			levenshtein.Lib(),
			fuzzy.Lib(),
			functional.Lib(),
			time.Lib(),
			urls.Lib(),
			fmt.Lib(),
		},
	}
}

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, map[string]lua.LValue{
		//"json":   json.New(L),
		//"html": html.New(L),
		//"crypto": crypto.New(L),
		//"http": http.New(L),
		//"regexp": regexp.New(L),
		//"js": js.New(L),
		//"cmd": cmd.New(L),
		//"strings": strings.New(L),
		//"pdf":  pdf.New(L),
		"time": time.New(L),
	}, nil)
}
