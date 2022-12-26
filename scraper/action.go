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

func (a *Action) Call(media []*Media) error {
	table := a.scraper.state.NewTable()
	for _, m := range media {
		table.Append(m.Value())
	}

	return a.scraper.state.CallByParam(lua.P{
		Fn:      a.Handler,
		NRet:    0,
		Protect: true,
	}, table, a.scraper.context)
}

func (s *Scraper) newAction(table *lua.LTable) (*Action, error) {
	action := &Action{scraper: s}
	err := gluamapper.Map(table, action)
	if err != nil {
		return nil, err
	}

	return action, nil
}
