package urls

import (
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/lua"
	"net/url"
)

const (
	valuesTypeName = "url_values"
	urlTypeName    = "url_url"
)

func Lib() *luadoc.Lib {
	classValues := &luadoc.Class{
		Name:        valuesTypeName,
		Description: "Values maps a string key to a list of values. It is typically used for query parameters and form values. Unlike in the `http.header` map, the keys in a `values` map are case-sensitive.",
		Methods: []*luadoc.Method{
			{
				Name:        "add",
				Description: "Adds the key and value to the values. It appends to any existing values associated with key.",
				Value:       urlValuesAdd,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key to add. It must not be empty.",
						Type:        luadoc.String,
					},
					{
						Name:        "value",
						Description: "The value to add. It must not be empty.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "set",
				Description: "Sets the key to value. It replaces any existing values associated with key.",
				Value:       urlValuesSet,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key to add. It must not be empty.",
						Type:        luadoc.String,
					},
					{
						Name:        "value",
						Description: "The value to add. It must not be empty.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "get",
				Description: "Gets the first value associated with the given key. If there are no values associated with the key, Get returns \"\".",
				Value:       urlValuesGet,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key to get.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: "The first value associated with the given key.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "has",
				Description: "Returns true if the values contains the specified key, false otherwise.",
				Value:       urlValuesHas,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key to check.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "has",
						Description: "True if the values contains the specified key, false otherwise.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "del",
				Description: "Deletes the values associated with key.",
				Value:       urlValuesDel,
				Params: []*luadoc.Param{
					{
						Name:        "key",
						Description: "The key to delete.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "encode",
				Description: "Encodes the values into \"URL encoded\" form sorted by key.",
				Value:       urlValuesEncode,
				Returns: []*luadoc.Param{
					{
						Name:        "encoded",
						Description: "The URL encoded form of the values.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "decode",
				Description: "Creates a values from the URL encoded form. It is the inverse operation of encode.",
				Value:       urlValuesDecode,
				Params: []*luadoc.Param{
					{
						Name:        "encoded",
						Description: "The URL encoded form of the values.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "values",
						Description: "The values created from the URL encoded form.",
						Type:        valuesTypeName,
					},
				},
			},
		},
	}

	classURL := &luadoc.Class{
		Name:        urlTypeName,
		Description: "Structured URL",
		Methods: []*luadoc.Method{
			{
				Name:        "hostname",
				Description: "",
				Value:       urlURLHostname,
				Returns: []*luadoc.Param{
					{
						Name:        "hostname",
						Description: "URLs hostname",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "parse",
				Description: "",
				Value:       urlURLParse,
				Params: []*luadoc.Param{
					{
						Name:        "ref",
						Description: "",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "url",
						Description: "",
						Type:        urlTypeName,
					},
				},
			},
			{
				Name:        "string",
				Description: "",
				Value:       urlURLString,
				Returns: []*luadoc.Param{
					{
						Name:        "url",
						Description: "",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "join_path",
				Description: "",
				Value:       urlURLJoinPath,
				Params: []*luadoc.Param{
					{
						Name:        "...",
						Description: "",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "url",
						Description: "",
						Type:        urlTypeName,
					},
				},
			},
			{
				Name:        "query",
				Description: "",
				Value:       urlURLQuery,
				Params: []*luadoc.Param{
					{
						Name:        "query",
						Description: "",
						Type:        valuesTypeName,
						Optional:    true,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "query",
						Description: "",
						Type:        valuesTypeName,
						Optional:    true,
					},
				},
			},
		},
	}

	return &luadoc.Lib{
		Name:        "urls",
		Description: "URLs is a library for working with URLs.",
		Funcs: []*luadoc.Func{
			{
				Name:        "parse",
				Description: "Parses URL",
				Value:       urlParse,
				Params: []*luadoc.Param{
					{
						Name:        "raw_url",
						Description: "URL string to parse",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "url",
						Description: "Parsed URL",
						Type:        urlTypeName,
					},
				},
			},
			{
				Name:        "values",
				Description: "Creates a new values.",
				Value:       urlValues,
				Returns: []*luadoc.Param{
					{
						Name:        "values",
						Description: "The new values.",
						Type:        valuesTypeName,
					},
				},
			},
			{
				Name:        "path_escape",
				Description: `Escapes the string so it can be safely placed inside a URL path segment, replacing special characters (including /) with %XX sequences as needed.`,
				Value:       urlPathEscape,
				Params: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The path to escape.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "escaped",
						Description: "The escaped path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "path_unescape",
				Description: `Unescapes a string; the inverse operation of path_escape. It converts each 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB.`,
				Value:       urlPathUnescape,
				Params: []*luadoc.Param{
					{
						Name:        "escaped",
						Description: "The escaped path.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The unescaped path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "query_escape",
				Description: `Escapes the string so it can be safely placed inside a URL query parameter, replacing special characters (including /) with %XX sequences as needed.`,
				Value:       urlQueryEscape,
				Params: []*luadoc.Param{
					{
						Name:        "query",
						Description: "The query to escape.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "escaped",
						Description: "The escaped query.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "query_unescape",
				Description: `Unescapes a string; the inverse operation of query_escape. It converts each 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB.`,
				Value:       urlQueryUnescape,
				Params: []*luadoc.Param{
					{
						Name:        "escaped",
						Description: "The escaped query.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "query",
						Description: "The unescaped query.",
						Type:        luadoc.String,
					},
					{
						Name:        "err",
						Description: "An error if the query could not be unescaped.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
		},
		Classes: []*luadoc.Class{
			classValues,
			classURL,
		},
	}
}

func checkValues(L *lua.LState, idx int) *url.Values {
	u, ok := L.CheckUserData(idx).Value.(*url.Values)
	if !ok {
		L.ArgError(idx, "expected `values`")
	}
	return u
}

func pushValues(L *lua.LState, v *url.Values) {
	ud := L.NewUserData()
	ud.Value = v
	L.SetMetatable(ud, L.GetTypeMetatable(valuesTypeName))
	L.Push(ud)
}

func checkURL(L *lua.LState, idx int) *url.URL {
	u, ok := L.CheckUserData(idx).Value.(*url.URL)
	if !ok {
		L.ArgError(idx, "expected `url`")
	}
	return u
}

func pushURL(L *lua.LState, u *url.URL) {
	ud := L.NewUserData()
	ud.Value = u
	L.SetMetatable(ud, L.GetTypeMetatable(urlTypeName))
	L.Push(ud)
}

func urlPathEscape(L *lua.LState) int {
	s := L.CheckString(1)
	L.Push(lua.LString(url.PathEscape(s)))
	return 1
}

func urlPathUnescape(L *lua.LState) int {
	s := L.CheckString(1)
	s, err := url.PathUnescape(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(s))
	return 1
}

func urlQueryEscape(L *lua.LState) int {
	s := L.CheckString(1)
	L.Push(lua.LString(url.QueryEscape(s)))
	return 1
}

func urlQueryUnescape(L *lua.LState) int {
	s := L.CheckString(1)
	s, err := url.QueryUnescape(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(s))
	return 1
}

func urlValues(L *lua.LState) int {
	pushValues(L, &url.Values{})
	return 1
}

func urlValuesAdd(L *lua.LState) int {
	v := checkValues(L, 1)
	key := L.CheckString(2)
	value := L.CheckString(3)
	v.Add(key, value)
	return 0
}

func urlValuesSet(L *lua.LState) int {
	v := checkValues(L, 1)
	key := L.CheckString(2)
	value := L.CheckString(3)
	v.Set(key, value)
	return 0
}

func urlValuesGet(L *lua.LState) int {
	v := checkValues(L, 1)
	key := L.CheckString(2)
	L.Push(lua.LString(v.Get(key)))
	return 1
}

func urlValuesEncode(L *lua.LState) int {
	v := checkValues(L, 1)
	L.Push(lua.LString(v.Encode()))
	return 1
}

func urlValuesDecode(L *lua.LState) int {
	s := L.CheckString(1)
	v, err := url.ParseQuery(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	pushValues(L, &v)
	return 1
}

func urlValuesDel(L *lua.LState) int {
	v := checkValues(L, 1)
	key := L.CheckString(2)
	v.Del(key)
	return 0
}

func urlValuesHas(L *lua.LState) int {
	v := checkValues(L, 1)
	key := L.CheckString(2)
	L.Push(lua.LBool(v.Has(key)))
	return 1
}

func urlParse(L *lua.LState) int {
	rawURL := L.CheckString(1)

	u, err := url.Parse(rawURL)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushURL(L, u)
	return 1
}

func urlURLHostname(L *lua.LState) int {
	u := checkURL(L, 1)

	L.Push(lua.LString(u.Hostname()))
	return 1
}

func urlURLQuery(L *lua.LState) int {
	u := checkURL(L, 1)

	if L.GetTop() == 1 {
		values := u.Query()
		pushValues(L, &values)
		return 1
	} else { // >= 2
		values := checkValues(L, 2)
		u.RawQuery = values.Encode()
		return 0
	}
}

func urlURLJoinPath(L *lua.LState) int {
	u := checkURL(L, 1)

	top := L.GetTop()
	elems := make([]string, top-1)

	for i := 2; i <= top; i++ {
		elems[i-2] = L.CheckString(i)
	}

	pushURL(L, u.JoinPath(elems...))
	return 1
}

func urlURLParse(L *lua.LState) int {
	u := checkURL(L, 1)
	ref := L.CheckString(2)

	parsed, err := u.Parse(ref)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	pushURL(L, parsed)
	return 1
}

func urlURLString(L *lua.LState) int {
	u := checkURL(L, 1)

	L.Push(lua.LString(u.String()))
	return 1
}
