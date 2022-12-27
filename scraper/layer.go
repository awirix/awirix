package scraper

import (
	"github.com/pkg/errors"
	"github.com/vivi-app/lua"
)

type Layer struct {
	scraper *Scraper
	*layer
}

type layer struct {
	Title   string
	Handler *lua.LFunction
	Noun    Noun `lua:"noun"`
}

func (l *Layer) String() string {
	if l.Title != "" {
		return l.Title
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
	}, value, l.scraper.context)
	if err != nil {
		return nil, err
	}

	return l.scraper.checkMediaSlice()
}

func (s *Scraper) newLayer(table *lua.LTable) (*Layer, error) {
	aux := &layer{}
	err := tableMapper.Map(table, aux)
	if err != nil {
		return nil, err
	}

	if aux.Title == "" {
		return nil, ErrMissingTitle
	}

	if aux.Handler == nil {
		return nil, errors.Wrap(ErrMissingHandler, "layer")
	}

	return &Layer{scraper: s, layer: aux}, nil
}
