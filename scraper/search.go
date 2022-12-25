package scraper

import (
	"github.com/vivi-app/gluamapper"
	"github.com/vivi-app/lua"
)

type Search struct {
	scraper     *Scraper
	Title       string
	Subtitle    string
	Placeholder string
	Handler     *lua.LFunction
	Noun        Noun
}

func (s *Search) String() string {
	if title := s.Title; title != "" {
		return title
	}

	return "Search"
}

func (s *Search) SubtitleAuto() string {
	if subtitle := s.Subtitle; subtitle != "" {
		return subtitle
	}

	return "Select a " + s.Noun.Singular()
}

func (s *Search) PlaceholderAuto() string {
	if placeholder := s.Placeholder; placeholder != "" {
		return placeholder
	}

	return "Search " + s.Noun.Plural()
}

func (s *Search) Call(query string) (subMedia []*Media, err error) {
	err = s.scraper.state.CallByParam(lua.P{
		Fn:      s.Handler,
		NRet:    1,
		Protect: true,
	}, lua.LString(query), s.scraper.progress)
	if err != nil {
		return nil, err
	}

	return s.scraper.checkMediaSlice()
}

func (s *Scraper) newSearch(table *lua.LTable) (*Search, error) {
	search := &Search{scraper: s}
	err := gluamapper.Map(table, search)
	if err != nil {
		return nil, err
	}

	return search, nil
}
