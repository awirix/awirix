package crypto

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/lualib/sdk/crypto/aes"
	"github.com/vivi-app/vivi/lualib/sdk/crypto/base64"
	"github.com/vivi-app/vivi/lualib/sdk/crypto/md5"
	"github.com/vivi-app/vivi/lualib/sdk/crypto/sha1"
	"github.com/vivi-app/vivi/lualib/sdk/crypto/sha256"
	"github.com/vivi-app/vivi/luautil"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, map[string]lua.LValue{
		"base64": base64.New(L),
		"md5":    md5.New(L),
		"sha1":   sha1.New(L),
		"sha256": sha256.New(L),
		"aes":    aes.New(L),
	}, nil)
}
