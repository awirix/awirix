package scraper

import (
	"github.com/vivi-app/gluamapper"
	"github.com/vivi-app/lua"
)

type Layer struct {
	scraper *Scraper
	Title   string
	Handler *lua.LFunction
	Noun    Noun
}

func (l *Layer) String() string {
	if title := l.Title; title != "" {
		return title
	}

	return "Select a " + l.Noun.Singular()
}

func (l *Layer) Call(media *Media) (subMedia []*Media, err error) {
	var value lua.LValue
	if media != nil {
		value = media.Value()
	} else {
		value = lua.LNil
	}

	err = l.scraper.state.CallByParam(lua.P{
		Fn:      l.Handler,
		NRet:    1,
		Protect: true,
	}, value, l.scraper.progress)
	if err != nil {
		return nil, err
	}

	return l.scraper.checkMediaSlice()
}

func (s *Scraper) newLayer(table *lua.LTable) (*Layer, error) {
	layer := &Layer{scraper: s}

	err := gluamapper.Map(table, layer)
	if err != nil {
		return nil, err
	}

	return layer, nil
}
