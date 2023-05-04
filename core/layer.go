package core

import (
	"github.com/awirix/lua"
	"github.com/pkg/errors"
)

type Layer struct {
	core  *Core
	cache map[*Media][]*Media
	*layer
}

type layer struct {
	Title   string
	Handler *lua.LFunction
	Noun    Noun `lua:"noun"`
}

func errLayer(err error) error {
	return errors.Wrap(err, "layer")
}

func (l *Layer) String() string {
	if l.Title != "" {
		return l.Title
	}

	return "Select a " + l.Noun.Singular()
}

func (l *Layer) Call(media *Media) (subMedia []*Media, err error) {
	if cached, ok := l.cache[media]; ok {
		return cached, nil
	}

	var value lua.LValue
	if media != nil {
		value = media.Value()
	} else {
		value = lua.LNil
	}

	err = l.core.state.CallByParam(lua.P{
		Fn:      l.Handler,
		NRet:    1,
		Protect: true,
	}, value, l.core.context)

	if err != nil {
		return nil, errLayer(err)
	}

	medias, err := l.core.checkMediaSlice()
	if err != nil {
		return nil, errLayer(err)
	}

	l.cache[media] = medias
	return medias, nil
}

func (c *Core) newLayer(table *lua.LTable) (*Layer, error) {
	aux := &layer{}
	err := tableMapper.Map(table, aux)

	if err != nil {
		return nil, errLayer(err)
	}

	if aux.Title == "" {
		return nil, errLayer(ErrMissingTitle)
	}

	if aux.Handler == nil {
		return nil, errLayer(ErrMissingHandler)
	}

	return &Layer{core: c, layer: aux, cache: make(map[*Media][]*Media)}, nil
}
