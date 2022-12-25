package scraper

import (
	"github.com/vivi-app/gluamapper"
	"github.com/vivi-app/lua"
)

type Action struct {
	scraper     *Scraper
	Title       string
	Description string
	Handler     *lua.LFunction
	Multiple    bool
}

func (a *Action) String() string {
	return a.Title
}

func (a *Action) Call(media *Media) error {
	return a.scraper.state.CallByParam(lua.P{
		Fn:      a.Handler,
		NRet:    1,
		Protect: true,
	}, media.Value(), a.scraper.progress)
}

func (s *Scraper) newAction(table *lua.LTable) (*Action, error) {
	action := &Action{scraper: s}
	err := gluamapper.Map(table, action)
	if err != nil {
		return nil, err
	}

	return action, nil
}
