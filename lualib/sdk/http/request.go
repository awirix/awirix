package http

import (
	lua "github.com/yuin/gopher-lua"
	"io"
	"net/http"
	"strings"
)

const requestTypeName = "request"

var requestMethods = map[string]lua.LGFunction{
	"set_header":     requestSetHeader,
	"set_basic_auth": requestSetBasicAuth,
}

func registerRequestType(L *lua.LState) {
	mt := L.NewTypeMetatable(requestTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), requestMethods))
}

func checkRequest(L *lua.LState, idx int) *http.Request {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*http.Request); ok {
		return v
	}
	L.ArgError(1, "request expected")
	return nil
}

func pushRequest(L *lua.LState, request *http.Request) {
	ud := L.NewUserData()
	ud.Value = request
	L.SetMetatable(ud, L.GetTypeMetatable(requestTypeName))
	L.Push(ud)
}

func newRequest(L *lua.LState) int {
	method := L.CheckString(1)
	url := L.CheckString(2)
	body := L.Get(3)

	var reqBody io.Reader

	if body.Type() == lua.LTNil {
		reqBody = nil
	} else if body.Type() == lua.LTString {
		reqBody = strings.NewReader(body.String())
	} else {
		L.ArgError(3, "string or nil expected")
		return 0
	}

	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushRequest(L, request)
	return 1
}

func requestSetHeader(L *lua.LState) int {
	request := checkRequest(L, 1)
	header := checkHeader(L, 2)
	request.Header = *header
	return 0
}

func requestSetBasicAuth(L *lua.LState) int {
	request := checkRequest(L, 1)
	username := L.CheckString(2)
	password := L.CheckString(3)
	request.SetBasicAuth(username, password)
	return 0
}

func requestDo(L *lua.LState) int {
	request := checkRequest(L, 1)
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushResponse(L, response)
	return 1
}
