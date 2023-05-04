package core

import (
	"github.com/awirix/lua"
	"github.com/pkg/errors"
)

type Action struct {
	core *Core
	*action
}

type action struct {
	Title       string
	Description string
	Handler     *lua.LFunction
}

func errAction(err error) error {
	return errors.Wrap(err, "action")
}

func (a *Action) String() string {
	return a.Title
}

func (a *Action) Call(media *Media) error {
	err := a.core.state.CallByParam(lua.P{
		Fn:      a.Handler,
		NRet:    0,
		Protect: true,
	}, media.Value(), a.core.context)

	if err != nil {
		return errAction(err)
	}

	return nil
}

func (c *Core) newAction(table *lua.LTable) (*Action, error) {
	aux := &action{}
	err := tableMapper.Map(table, aux)

	if err != nil {
		return nil, errAction(err)
	}

	return &Action{core: c, action: aux}, nil
}
