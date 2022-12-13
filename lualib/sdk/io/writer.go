package io

import (
	lua "github.com/vivi-app/lua"
	"io"
)

const writerTypeName = "writer"

func registerWriterType(L *lua.LState) {
	mt := L.NewTypeMetatable(writerTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), writerMethods))
}

func PushWriter(L *lua.LState, w io.Writer) {
	ud := L.NewUserData()
	ud.Value = w
	L.SetMetatable(ud, L.GetTypeMetatable(writerTypeName))
	L.Push(ud)
}

var writerMethods = map[string]lua.LGFunction{
	"write": writerWrite,
}

func checkWriter(L *lua.LState, n int) io.Writer {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(io.Writer); ok {
		return v
	}
	L.ArgError(1, "writer expected")
	return nil
}

func writerWrite(L *lua.LState) int {
	w := checkWriter(L, 1)
	data := L.CheckString(2)
	n, err := w.Write([]byte(data))
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(n))
	return 1
}

func newNopWriter(L *lua.LState) int {
	PushWriter(L, &nopWriter{})
	return 1
}

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
