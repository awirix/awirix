package http

import (
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

const responseTypeName = "response"

var responseMethods = map[string]lua.LGFunction{
	"status_code":       responseStatusCode,
	"body":              responseBody,
	"headers":           responseHeaders,
	"cookies":           responseCookies,
	"status":            responseStatus,
	"content_length":    responseContentLength,
	"transfer_encoding": responseTransferEncoding,
	"proto":             responseProto,
	"proto_major":       responseProtoMajor,
	"proto_minor":       responseProtoMinor,
	"trailer":           responseTrailer,
	"close":             responseClose,
	"uncompressed":      responseUncompressed,
	"request":           responseRequest,
}

func registerResponseType(L *lua.LState) {
	mt := L.NewTypeMetatable(responseTypeName)
	L.SetGlobal("response", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), responseMethods))
}

func pushResponse(L *lua.LState, response *http.Response) {
	ud := L.NewUserData()
	ud.Value = response
	L.SetMetatable(ud, L.GetTypeMetatable(responseTypeName))
	L.Push(ud)
}

func checkResponse(L *lua.LState, n int) *http.Response {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*http.Response); ok {
		return v
	}
	L.ArgError(n, "response expected")
	return nil
}

func responseStatusCode(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LNumber(response.StatusCode))
	return 1
}

func responseBody(L *lua.LState) int {
	response := checkResponse(L, 1)

	var b []byte
	_, err := response.Body.Read(b)
	defer response.Body.Close()

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(b))
	return 1
}

func responseHeaders(L *lua.LState) int {
	response := checkResponse(L, 1)
	headers := response.Header

	pushHeader(L, &headers)
	return 1
}

func responseCookies(L *lua.LState) int {
	response := checkResponse(L, 1)
	cookies := response.Cookies()

	table := L.NewTable()
	for i, cookie := range cookies {
		pushCookie(L, cookie)
		table.RawSetInt(i, L.Get(-1))
		L.Pop(1)
	}

	L.Push(table)
	return 1
}

func responseStatus(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LString(response.Status))
	return 1
}

func responseContentLength(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LNumber(response.ContentLength))
	return 1
}

func responseTransferEncoding(L *lua.LState) int {
	response := checkResponse(L, 1)
	transferEncoding := response.TransferEncoding

	table := L.NewTable()
	for i, encoding := range transferEncoding {
		table.RawSetInt(i, lua.LString(encoding))
	}

	L.Push(table)
	return 1
}

func responseProto(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LString(response.Proto))
	return 1
}

func responseProtoMajor(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LNumber(response.ProtoMajor))
	return 1
}

func responseProtoMinor(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LNumber(response.ProtoMinor))
	return 1
}

func responseTrailer(L *lua.LState) int {
	response := checkResponse(L, 1)
	trailer := response.Trailer

	pushHeader(L, &trailer)
	return 1
}

func responseClose(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LBool(response.Close))
	return 1
}

func responseUncompressed(L *lua.LState) int {
	response := checkResponse(L, 1)
	L.Push(lua.LBool(response.Uncompressed))
	return 1
}

func responseRequest(L *lua.LState) int {
	response := checkResponse(L, 1)
	pushRequest(L, response.Request)
	return 1
}
