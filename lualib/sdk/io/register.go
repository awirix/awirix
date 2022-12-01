package io

import (
	"github.com/vivi-app/vivi/luautil"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	registerReaderCloserType(L)
	registerWriterType(L)

	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"reader_closer_from_string": readerCloserFromString,
		"nop_writer":                newNopWriter,
	})
}
