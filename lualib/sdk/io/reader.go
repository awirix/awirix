package io

import (
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/key"
	lua "github.com/yuin/gopher-lua"
	"io"
	"strings"
)

const readCloserTypeName = "read_closer"

func registerReaderCloserType(L *lua.LState) {
	mt := L.NewTypeMetatable(readCloserTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), readCloserMethods))
}

func PushReadCloser(L *lua.LState, r io.ReadCloser) {
	ud := L.NewUserData()
	ud.Value = r
	L.SetMetatable(ud, L.GetTypeMetatable(readCloserTypeName))
	L.Push(ud)
}

var readCloserMethods = map[string]lua.LGFunction{
	"read":  readCloserRead,
	"close": readCloserClose,
}

func checkReadCloser(L *lua.LState, n int) io.ReadCloser {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(io.ReadCloser); ok {
		return v
	}
	L.ArgError(1, "reader_closer expected")
	return nil
}

func readCloserRead(L *lua.LState) int {
	r := checkReadCloser(L, 1)
	data, err := io.ReadAll(r)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer func() {
		if viper.GetBool(key.ExtensionsIOAutoCloseReaders) {
			_ = r.Close()
		}
	}()

	L.Push(lua.LString(data))
	return 1
}

func readerCloserFromString(L *lua.LState) int {
	s := L.CheckString(1)
	PushReadCloser(L, io.NopCloser(strings.NewReader(s)))
	return 1
}

func readCloserClose(L *lua.LState) int {
	r := checkReadCloser(L, 1)
	err := r.Close()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
