package scraper

import (
	"github.com/pkg/errors"
	"github.com/vivi-app/lua"
)

type Action struct {
	scraper *Scraper
	*action
}

type action struct {
	Title       string
	Description string
	Handler     *lua.LFunction
	Max         int
}

func (a *Action) String() string {
	return a.Title
}

func (a *Action) Call(media []*Media) error {
	if a.Max != 0 && len(media) > a.Max {
		return errors.New("too many medias")
	}

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
	aux := &action{}
	err := tableMapper.Map(table, aux)
	if err != nil {
		return nil, errors.Wrap(err, "action")
	}

	return &Action{scraper: s, action: aux}, nil
}
