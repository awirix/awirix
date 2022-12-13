package io

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	registerReaderCloserType(L)
	registerWriterType(L)

	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"reader_closer_from_string": readerCloserFromString,
		"nop_writer":                newNopWriter,
	})
}
