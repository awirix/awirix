package scraper

import "github.com/vivi-app/lua"

type Context struct {
	Progress func(message string)
	Error    func(err error)

	progress *lua.LFunction
	error    *lua.LFunction
}
