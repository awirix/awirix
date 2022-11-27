package aes

import (
	"crypto/aes"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"encrypt": encrypt,
		"decrypt": decrypt,
	})
}

func encrypt(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)

	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	encrypted := make([]byte, len(value))
	cipher.Encrypt(encrypted, []byte(value))

	L.Push(lua.LString(encrypted))
	return 1
}

func decrypt(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)

	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	decrypted := make([]byte, len(value))
	cipher.Decrypt(decrypted, []byte(value))

	L.Push(lua.LString(decrypted))
	return 1
}
