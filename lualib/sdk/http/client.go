package http

import (
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

const clientTypeName = "client"

func registerClientType(L *lua.LState) {
	mt := L.NewTypeMetatable(clientTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), clientMethods))
}

func pushClient(L *lua.LState, client *http.Client) {
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable(clientTypeName))
	L.Push(ud)
}

func checkClient(L *lua.LState, n int) *http.Client {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*http.Client); ok {
		return v
	}
	L.ArgError(n, "client expected")
	return nil
}

var clientMethods = map[string]lua.LGFunction{
	"send": clientDo,
}

func newClient(L *lua.LState) int {
	client := &http.Client{}
	pushClient(L, client)
	return 1
}

func clientDo(L *lua.LState) int {
	client := checkClient(L, 1)
	req := checkRequest(L, 2)

	resp, err := client.Do(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushResponse(L, resp)
	return 1
}
