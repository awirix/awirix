package sha512

import (
	"crypto/sha512"
	"github.com/awirix/awirix/luadoc"
	lua "github.com/awirix/lua"
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

func sum(L *lua.LState) int {
	value := L.CheckString(1)
	s := sha512.Sum512([]byte(value))
	L.Push(lua.LString(s[:]))
	return 1
}
