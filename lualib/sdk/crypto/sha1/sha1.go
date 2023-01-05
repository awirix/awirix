package sha1

import (
	"crypto/sha1"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luadoc"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "sha1",
		Description: "SHA1 cryptographic hash function.",
		Funcs: []*luadoc.Func{
			{
				Name:        "sum",
				Description: "Returns the SHA1 hash of the given string.",
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
						Description: "The SHA1 hash of the given string.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func sum(L *lua.LState) int {
	value := L.CheckString(1)
	s := sha1.Sum([]byte(value))
	L.Push(lua.LString(s[:]))
	return 1
}
