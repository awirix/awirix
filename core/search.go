package core

import (
	"github.com/awirix/lua"
	"github.com/pkg/errors"
	"strings"
)

type Search struct {
	core  *Core
	cache map[string][]*Media
	*search
}

type search struct {
	Title       string
	Subtitle    string
	Placeholder string
	Handler     *lua.LFunction
	Noun        Noun `lua:"noun"`
}

func errSearch(err error) error {
	return errors.Wrap(err, "search")
}

func (s *Search) String() string {
	if s.Title != "" {
		return s.Title
	}

	return "Search"
}

func (s *Search) Placeholder() string {
	if s.search.Placeholder != "" {
		return s.search.Placeholder
	}

	return "Search for " + s.search.Noun.Plural()
}

func (s *Search) Subtitle() string {
	if s.search.Subtitle != "" {
		return s.search.Subtitle
	}

	return "Select a " + s.search.Noun.Singular()
}

func (s *Search) Call(query string) (subMedia []*Media, err error) {
	query = strings.TrimSpace(query)

	if cached, ok := s.cache[query]; ok {
		return cached, nil
	}

	err = s.core.state.CallByParam(lua.P{
		Fn:      s.Handler,
		NRet:    1,
		Protect: true,
	}, lua.LString(query), s.core.context)

	if err != nil {
		return nil, errSearch(err)
	}

	media, err := s.core.checkMediaSlice()
	if err != nil {
		return nil, errSearch(err)
	}

	s.cache[query] = media
	return media, nil
}

func (c *Core) newSearch(table *lua.LTable) (*Search, error) {
	aux := &search{}
	err := tableMapper.Map(table, aux)

	if err != nil {
		return nil, errSearch(err)
	}

	if aux.Handler == nil {
		return nil, errSearch(ErrMissingHandler)
	}

	return &Search{core: c, search: aux, cache: make(map[string][]*Media)}, nil
}
