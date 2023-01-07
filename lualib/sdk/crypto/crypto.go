package crypto

import (
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/awirix/lualib/sdk/crypto/aes"
	"github.com/awirix/awirix/lualib/sdk/crypto/base64"
	"github.com/awirix/awirix/lualib/sdk/crypto/md5"
	"github.com/awirix/awirix/lualib/sdk/crypto/sha1"
	"github.com/awirix/awirix/lualib/sdk/crypto/sha256"
	"github.com/awirix/awirix/lualib/sdk/crypto/sha512"
	lua "github.com/awirix/lua"
)

func Lib(L *lua.LState) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "crypto",
		Description: "Various cryptographic functions.",
		Libs: []*luadoc.Lib{
			base64.Lib(L),
			aes.Lib(),
			md5.Lib(),
			sha1.Lib(),
			sha256.Lib(),
			sha512.Lib(),
		},
	}
}
