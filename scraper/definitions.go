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

type Layer struct {
	scraper  *Scraper
	Name     string
	function *lua.LFunction
}

func (l *Layer) String() string {
	return l.Name
}

func (l *Layer) Call(media *Media) (subMedias []*Media, err error) {
	var value lua.LValue
	if media != nil {
		value = media.Value()
	} else {
		value = lua.LNil
	}

	err = l.scraper.state.CallByParam(lua.P{
		Fn:      l.function,
		NRet:    1,
		Protect: true,
	}, value, l.scraper.progress)
	if err != nil {
		return nil, err
	}

	return l.scraper.checkMediaSlice()
}
