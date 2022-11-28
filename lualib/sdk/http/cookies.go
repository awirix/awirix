package http

import (
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

const cookieTypeName = "cookie"

func registerCookieType(L *lua.LState) {
	mt := L.NewTypeMetatable(cookieTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), cookieMethods))
}

func pushCookie(L *lua.LState, cookie *http.Cookie) {
	ud := L.NewUserData()
	ud.Value = cookie
	L.SetMetatable(ud, L.GetTypeMetatable(cookieTypeName))
	L.Push(ud)
}

func checkCookie(L *lua.LState, n int) *http.Cookie {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*http.Cookie); ok {
		return v
	}
	L.ArgError(n, "cookie expected")
	return nil
}

var cookieMethods = map[string]lua.LGFunction{
	"name":  cookieName,
	"value": cookieValue,
}

func cookieName(L *lua.LState) int {
	cookie := checkCookie(L, 1)
	L.Push(lua.LString(cookie.Name))
	return 1
}

func cookieValue(L *lua.LState) int {
	cookie := checkCookie(L, 1)
	L.Push(lua.LString(cookie.Value))
	return 1
}
