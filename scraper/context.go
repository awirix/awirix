package scraper

import "github.com/awirix/lua"

type Context struct {
	Progress func(message string)
	Error    func(err error)

	progress *lua.LFunction
	error    *lua.LFunction
}
