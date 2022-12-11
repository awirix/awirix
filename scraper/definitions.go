package scraper

import lua "github.com/yuin/gopher-lua"

const Module = "scraper"

const (
	FunctionSearch   = "search"
	FunctionPrepare  = "prepare"
	FunctionStream   = "stream"
	FunctionDownload = "download"
)

const (
	FieldDisplay = "display"
	FieldAbout   = "about"
	FieldLayers  = "layers"
)

type Layer struct {
	Name        string
	Function    func(media *Media) (subMedias []*Media, err error)
	luaFunction *lua.LFunction
}