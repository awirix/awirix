package http

import (
	"github.com/awirix/awirix/luadoc"
	lua "github.com/awirix/lua"
	"net/http"
	"strings"
)

const httpTypeName = "http"

func Lib() *luadoc.Lib {
	httpMethods := luadoc.Enum(
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPut,
		http.MethodPatch,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
	)

	classCookie := &luadoc.Class{
		Name:        cookieTypeName,
		Description: "A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an HTTP response or the Cookie header of an HTTP request.",
		Methods: []*luadoc.Method{
			{
				Name:        "name",
				Description: "Returns the name of the cookie.",
				Value:       cookieName,
				Returns: []*luadoc.Param{
					{
						Name:        "name",
						Description: "The name of the cookie.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "value",
				Description: "Returns the value of the cookie.",
				Value:       cookieValue,
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: "The value of the cookie.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}

	classHeader := &luadoc.Class{
		Name:        headerTypeName,
		Description: "A Header represents the key-value pairs in an HTTP header.",
		Methods: []*luadoc.Method{
			{
				Name:        "add",
				Description: "Adds the key, value pair to the header. It appends to any existing values associated with key.",
				Value:       headerAdd,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key of the header.",
						Type:        luadoc.String,
					},
					{
						Name:        "value",
						Description: "The value of the header.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "set",
				Description: "Sets the header entries associated with key to the single element value. It replaces any existing values associated with key.",
				Value:       headerSet,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key of the header.",
						Type:        luadoc.String,
					},
					{
						Name:        "value",
						Description: "The value of the header.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "get",
				Description: "Gets the first value associated with the given key. If there are no values associated with the key, Get returns \"\".",
				Value:       headerGet,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key of the header.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: "The value of the header.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "del",
				Description: "Deletes the values associated with key.",
				Value:       headerDel,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key of the header.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "clone",
				Description: "Returns a copy of the header.",
				Value:       headerClone,
				Returns: []*luadoc.Param{
					{
						Name:        "header",
						Description: "A copy of the header.",
						Type:        headerTypeName,
					},
				},
			},
		},
	}

	classResponse := &luadoc.Class{
		Name:        responseTypeName,
		Description: "Represents an HTTP response.",
		Methods: []*luadoc.Method{
			{
				Name:        "status",
				Description: "Returns the HTTP status code of the response.",
				Value:       responseStatus,
				Returns: []*luadoc.Param{
					{
						Name:        "status",
						Description: "The HTTP status code of the response.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "status_code",
				Description: "Returns the HTTP status code of the response.",
				Value:       responseStatusCode,
				Returns: []*luadoc.Param{
					{
						Name:        "code",
						Description: "The HTTP status code of the response.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "body",
				Description: "Returns the body of the response.",
				Value:       responseBody,
				Returns: []*luadoc.Param{
					{
						Name:        "body",
						Description: "The body of the response.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "content_length",
				Description: "Returns the content length of the response.",
				Value:       responseContentLength,
				Returns: []*luadoc.Param{
					{
						Name:        "length",
						Description: "The content length of the response.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "transfer_encoding",
				Description: "Returns the transfer encoding of the response.",
				Value:       responseTransferEncoding,
				Returns: []*luadoc.Param{
					{
						Name:        "encoding",
						Description: "The transfer encoding of the response.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "header",
				Description: "Returns the header of the response.",
				Value:       responseHeader,
				Returns: []*luadoc.Param{
					{
						Name:        "header",
						Description: "The header of the response.",
						Type:        headerTypeName,
					},
				},
			},
			{
				Name:        "cookies",
				Description: "Returns the cookies of the response.",
				Value:       responseCookies,
				Returns: []*luadoc.Param{
					{
						Name:        "cookies",
						Description: "The cookies of the response.",
						Type:        luadoc.List(cookieTypeName),
					},
				},
			},
		},
	}

	classRequest := &luadoc.Class{
		Name:        requestTypeName,
		Description: "Represents an HTTP request.",
		Methods: []*luadoc.Method{
			{
				Name:        "set_header",
				Description: "Sets the header to the request.",
				Value:       requestSetHeader,
				Params: []*luadoc.Param{
					{
						Name:        "header",
						Description: "The header to set.",
						Type:        headerTypeName,
					},
				},
			},
			{
				Name:        "get_header",
				Description: "Gets the header of the request.",
				Value:       requestGetHeader,
				Returns: []*luadoc.Param{
					{
						Name:        "header",
						Description: "The header of the request.",
						Type:        headerTypeName,
					},
				},
			},
			{
				Name:        "set_basic_auth",
				Description: "Sets the basic authentication to the request.",
				Value:       requestSetBasicAuth,
				Params: []*luadoc.Param{
					{
						Name:        "username",
						Description: "The username of the basic authentication.",
						Type:        luadoc.String,
					},
					{
						Name:        "password",
						Description: "The password of the basic authentication.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "set_method",
				Description: "Sets the method to the request.",
				Value:       requestSetMethod,
				Params: []*luadoc.Param{
					{
						Name:        "method",
						Description: "The method to set.",
						Type:        httpMethods,
					},
				},
			},
			{
				Name:        "get_method",
				Description: "Gets the method of the request.",
				Value:       requestGetMethod,
				Returns: []*luadoc.Param{
					{
						Name:        "method",
						Description: "The method of the request.",
						Type:        httpMethods,
					},
				},
			},
			{
				Name:        "set_url",
				Description: "Sets the URL to the request.",
				Value:       requestSetURL,
				Params: []*luadoc.Param{
					{
						Name:        "url",
						Description: "The URL to set.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "url",
						Description: "The URL of the request.",
						Type:        luadoc.String,
					},
					{
						Name:        "err",
						Description: "An error if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "get_url",
				Description: "Gets the URL of the request.",
				Value:       requestGetURL,
				Returns: []*luadoc.Param{
					{
						Name:        "url",
						Description: "The URL of the request.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "set_body",
				Description: "Sets the body to the request.",
				Value:       requestSetBody,
				Params: []*luadoc.Param{
					{
						Name:        "body",
						Description: "The body to set.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "get_body",
				Description: "Gets the body of the request.",
				Value:       requestGetBody,
				Returns: []*luadoc.Param{
					{
						Name:        "body",
						Description: "The body of the request.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}

	classClient := &luadoc.Class{
		Name:        clientTypeName,
		Description: "Represents an HTTP client.",
		Methods: []*luadoc.Method{
			{
				Name:        "send",
				Description: "Sends the request and returns the response.",
				Value:       clientSend,
				Params: []*luadoc.Param{
					{
						Name:        "request",
						Description: "The request to send.",
						Type:        requestTypeName,
					},
					{
						Name:        "cache",
						Description: "Whether to cache the response.",
						Type:        luadoc.Boolean,
						Optional:    true,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "response",
						Description: "The response of the request.",
						Type:        responseTypeName,
					},
					{
						Name:        "error",
						Description: "The error if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "set_timeout",
				Description: "Sets the timeout of the client.",
				Value:       clientSetTimeout,
				Params: []*luadoc.Param{
					{
						Name:        "timeout",
						Description: "The timeout of the client in seconds.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "batch",
				Description: "Sends the requests in batch and returns the responses.",
				Value:       clientSendBatch,
				Params: []*luadoc.Param{
					{
						Name:        "requests",
						Description: "The requests to send.",
						Type:        luadoc.Map(luadoc.String, requestTypeName),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "responses",
						Description: "The responses of the requests.",
						Type:        luadoc.Map(luadoc.String, responseTypeName),
					},
					{
						Name:        "error",
						Description: "The error if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
		},
	}

	return &luadoc.Lib{
		Name:        httpTypeName,
		Description: "HTTP client library.",
		Funcs: []*luadoc.Func{
			{
				Name:        "get",
				Description: "Performs a GET request to the specified URL.",
				Value:       defaultClientGet,
				Params: []*luadoc.Param{
					{
						Name:        "url",
						Description: "The URL to request.",
						Type:        luadoc.String,
					},
					{
						Name:        "cache",
						Description: "Whether to cache the response.",
						Type:        luadoc.Boolean,
						Optional:    true,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "response",
						Description: "The response object.",
						Type:        responseTypeName,
					},
					{
						Name:        "error",
						Description: "The error message, if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "post",
				Description: "Performs a POST request to the specified URL.",
				Value:       defaultClientPost,
				Params: []*luadoc.Param{
					{
						Name:        "url",
						Description: "The URL to request.",
						Type:        luadoc.String,
					},
					{
						Name:        "body",
						Description: "The body of the request.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "response",
						Description: "The response object.",
						Type:        responseTypeName,
					},
					{
						Name:        "error",
						Description: "The error message, if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "request",
				Description: "Creates a new request object.",
				Value:       newRequest,
				Params: []*luadoc.Param{
					{
						Name:        "method",
						Description: "The method of the request.",
						Type: luadoc.Enum(
							http.MethodGet,
							http.MethodPost,
							http.MethodDelete,
							http.MethodPut,
							http.MethodPatch,
							http.MethodConnect,
							http.MethodHead,
							http.MethodOptions,
						),
					},
					{
						Name:        "url",
						Description: "The URL of the request.",
						Type:        luadoc.String,
					},
					{
						Name:        "body",
						Description: "The body of the request.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "request",
						Description: "The request object.",
						Type:        requestTypeName,
					},
					{
						Name:        "error",
						Description: "The error message, if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "client",
				Description: "Creates a new client object.",
				Value:       newClient,
				Returns: []*luadoc.Param{
					{
						Name:        "client",
						Description: "The client object.",
						Type:        clientTypeName,
					},
				},
			},
			{
				Name:        "header",
				Description: "Creates a new header object.",
				Value:       newHeader,
				Returns: []*luadoc.Param{
					{
						Name:        "header",
						Description: "The header object.",
						Type:        headerTypeName,
					},
				},
			},
		},
		Classes: []*luadoc.Class{
			classCookie,
			classHeader,
			classResponse,
			classRequest,
			classClient,
		},
	}
}

func defaultClientGet(L *lua.LState) int {
	url := L.CheckString(1)
	doCache := L.OptBool(2, false)

	// error can not occur here
	req, _ := http.NewRequest("GET", url, nil)

	if res, ok := cacheGet(L, req); ok {
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
		_ = cacheSet(L, req, res)
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
