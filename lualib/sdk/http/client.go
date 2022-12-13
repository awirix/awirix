package http

import (
	"github.com/vivi-app/lua"
	"net/http"
	"sync"
	"time"
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
	"send":        clientSend,
	"batch":       clientSendBatch,
	"set_timeout": clientSetTimeout,
}

func newClient(L *lua.LState) int {
	client := &http.Client{}
	pushClient(L, client)
	return 1
}

func clientSend(L *lua.LState) int {
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

func clientSetTimeout(L *lua.LState) int {
	client := checkClient(L, 1)
	seconds := L.CheckNumber(2)

	client.Timeout = time.Second * time.Duration(seconds)
	return 0
}

func clientSendBatch(L *lua.LState) int {
	client := checkClient(L, 1)

	var requests = make(map[lua.LValue]*http.Request, 0)

	L.CheckTable(2).ForEach(func(key, value lua.LValue) {
		req, ok := value.(*lua.LUserData).Value.(*http.Request)
		if !ok {
			L.ArgError(2, "table of requests expected")
			return
		}

		requests[key] = req
	})

	var (
		responses = L.NewTable()
		err       error
	)

	var wg sync.WaitGroup
	wg.Add(len(requests))

	for key, req := range requests {
		go func(key lua.LValue, req *http.Request) {
			defer wg.Done()

			if err != nil {
				return
			}

			var resp *http.Response
			resp, err = client.Do(req)
			if err != nil {
				return
			}

			ud := L.NewUserData()
			ud.Value = resp
			L.SetMetatable(ud, L.GetTypeMetatable(responseTypeName))
			responses.RawSet(key, ud)
		}(key, req)
	}

	wg.Wait()

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(responses)
	return 1
}
