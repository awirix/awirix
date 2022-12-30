package http

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/cache"
	"github.com/vivi-app/vivi/luautil"
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
		"get":     defaultClientGet,
		"post":    defaultClientPost,
		"request": newRequest,
		"header":  newHeader,
		"client":  newClient,
	})
}

func defaultClientGet(L *lua.LState) int {
	url := L.CheckString(1)
	doCache := L.OptBool(3, false)

	// error can not occur here
	req, _ := http.NewRequest("GET", url, nil)

	if res, ok := cache.HTTP.Get(req); ok {
		pushResponse(L, res)
		return 1
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if doCache {
		_ = cache.HTTP.Set(req, res)
	}

	pushResponse(L, res)
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
