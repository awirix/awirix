package sha512

import (
	"crypto/sha512"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luadoc"
	"github.com/vivi-app/vivi/luautil"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "sha512",
		Description: "SHA-512 cryptographic hash function.",
		Funcs: []*luadoc.Func{
			{
				Name:        "sum",
				Description: "Returns the SHA-512 hash of the given string.",
				Value:       sum,
				Params: []*luadoc.Param{
					{
						Name:        "value",
						Description: "The string to hash.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "hash",
						Description: "The SHA-512 hash of the given string.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"sum": sum,
	})
}

func sum(L *lua.LState) int {
	value := L.CheckString(1)
	s := sha512.Sum512([]byte(value))
	L.Push(lua.LString(s[:]))
	return 1
}
