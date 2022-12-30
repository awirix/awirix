package scraper

import (
	"fmt"
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
	Min, Max    int
}

func errAction(err error) error {
	return errors.Wrap(err, "action")
}

func (a *Action) String() string {
	return a.Title
}

func (a *Action) InBounds(length int) bool {
	if a.Min != 0 {
		if length < a.Min {
			return false
		}
	}

	if a.Max != 0 {
		if length > a.Max {
			return false
		}
	}

	return true
}

func (a *Action) BoundsString() string {
	if a.Min == 0 && a.Max == 0 {
		return ""
	}

	if a.Min == 0 {
		return "max " + fmt.Sprint(a.Max)
	}

	if a.Max == 0 {
		return "min " + fmt.Sprint(a.Min)
	}

	return fmt.Sprint(a.Min) + ".." + fmt.Sprint(a.Max)
}

func (a *Action) Call(media []*Media) error {
	if !a.InBounds(len(media)) {
		return errAction(errors.New("out of bounds"))
	}

	table := a.scraper.state.NewTable()
	for _, m := range media {
		table.Append(m.Value())
	}

	err := a.scraper.state.CallByParam(lua.P{
		Fn:      a.Handler,
		NRet:    0,
		Protect: true,
	}, table, a.scraper.context)

	if err != nil {
		return errAction(err)
	}

	return nil
}

func (s *Scraper) newAction(table *lua.LTable) (*Action, error) {
	aux := &action{}
	err := tableMapper.Map(table, aux)

	if err != nil {
		return nil, errAction(err)
	}

	if aux.Min > aux.Max {
		return nil, errAction(errors.New("min is greater than max"))
	}

	return &Action{scraper: s, action: aux}, nil
}
