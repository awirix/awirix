package sha256

import (
	"crypto/sha256"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luadoc"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "sha256",
		Description: "SHA256 cryptographic hash function.",
		Funcs: []*luadoc.Func{
			{
				Name:        "sum",
				Description: "Returns the SHA256 hash of the given string.",
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
						Description: "The SHA256 hash of the given string.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func sum(L *lua.LState) int {
	value := L.CheckString(1)
	s := sha256.Sum256([]byte(value))
	L.Push(lua.LString(s[:]))
	return 1
}
