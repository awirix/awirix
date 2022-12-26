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
	Noun        *Noun
}

func (s *Search) String() string {
	return s.Title
}

func (s *Search) Call(query string) (subMedia []*Media, err error) {
	err = s.scraper.state.CallByParam(lua.P{
		Fn:      s.Handler,
		NRet:    1,
		Protect: true,
	}, lua.LString(query), s.scraper.context)
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

	if search.Noun == nil {
		search.Noun = &Noun{singular: "media"}
	}

	if search.Title == "" {
		search.Title = "Search"
	}

	if search.Subtitle == "" {
		search.Subtitle = "Select a " + search.Noun.Singular()
	}

	if search.Placeholder == "" {
		search.Placeholder = "Search for " + search.Noun.Plural()
	}

	return search, nil
}
