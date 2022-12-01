package http

import (
	"github.com/vivi-app/vivi/luautil"
	lua "github.com/yuin/gopher-lua"
	"net/http"
	"strings"
)

func New(L *lua.LState) *lua.LTable {
	registerResponseType(L)
	registerRequestType(L)
	registerHeaderType(L)
	registerClientType(L)
	registerCookieType(L)

	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"get":         defaultClientGet,
		"post":        defaultClientPost,
		"new_request": newRequest,
		"new_header":  newHeader,
		"new_client":  newClient,
	})
}

func defaultClientGet(L *lua.LState) int {
	url := L.CheckString(1)
	response, err := http.Get(url)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushResponse(L, response)
	return 1
}

func defaultClientPost(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.CheckString(2)
	response, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushResponse(L, response)
	return 1
}
