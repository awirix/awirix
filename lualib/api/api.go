package api

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"play_video": playVideo,
		"open":       openDefault,
		"download":   download,
	})
}
