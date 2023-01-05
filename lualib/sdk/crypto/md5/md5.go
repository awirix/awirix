package md5

import (
	"crypto/md5"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luadoc"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "md5",
		Description: "MD5 cryptographic hash function.",
		Funcs: []*luadoc.Func{
			{
				Name:        "sum",
				Description: "Returns the MD5 hash of the given string.",
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
						Description: "The MD5 hash of the given string.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func sum(L *lua.LState) int {
	value := L.CheckString(1)
	s := md5.Sum([]byte(value))
	L.Push(lua.LString(s[:]))
	return 1
}
