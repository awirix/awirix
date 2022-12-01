package api

import (
	"github.com/vivi-app/vivi/luautil"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"play_video": playVideo,
		"open":       openDefault,
		"download":   download,
	})
}
