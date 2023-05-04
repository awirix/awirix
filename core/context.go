package core

import (
	"github.com/awirix/lua"
	"github.com/samber/mo"
)

type Context struct {
	Progress func(message string)
	Error    func(err error)

	query    mo.Option[string]
	selected []*Media

	progress *lua.LFunction
	error    *lua.LFunction
}
