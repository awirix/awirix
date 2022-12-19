package scraper

import lua "github.com/vivi-app/lua"

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

type LayerTakesType int

const (
	LayerTakesNone LayerTakesType = iota + 1
	LayerTakesString
	LayerTakesMedia
)

type Layer struct {
	Name        string
	Takes       LayerTakesType
	Function    func(media *Media) (subMedias []*Media, err error)
	luaFunction *lua.LFunction
}
