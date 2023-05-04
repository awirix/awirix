package core

import (
	"fmt"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/luautil"
	"github.com/awirix/lua"
	"github.com/pkg/errors"
	"io"
)

type Core struct {
	state *lua.LState

	search  *Search
	layers  []*Layer
	actions []*Action
	context *lua.LTable
}

func errCore(err error) error {
	return errors.Wrap(err, "core")
}

func (c *Core) HasSearch() bool {
	return c.search != nil
}

func (c *Core) HasLayers() bool {
	return c.layers != nil && len(c.layers) > 0
}

func (c *Core) HasActions() bool {
	return c.actions != nil && len(c.actions) > 0
}

func (c *Core) SetExtensionContext(context *Context) {
	c.context = luautil.NewTable(c.state, nil, map[string]lua.LGFunction{
		"progress": func(L *lua.LState) int {
			message := L.ToString(1)
			context.Progress(message)
			log.Tracef("progress: %c", message)
			return 0
		},
		"error": func(L *lua.LState) int {
			err := errors.New(L.ToString(1))
			context.Error(err)
			log.Tracef("error: %c", err)
			return 0
		},
	})
}

func (c *Core) getLayers(table *lua.LTable) (layers []*Layer, err error) {
	field := table.RawGetString(FieldLayers)

	if field.Type() == lua.LTNil {
		return nil, nil
	} else if field.Type() != lua.LTTable {
		return nil, fmt.Errorf("layers must be a table, got %c", field.Type().String())
	}

	table = field.(*lua.LTable)

	table.ForEach(func(_, value lua.LValue) {
		if err != nil {
			return
		}

		if value.Type() != lua.LTTable {
			err = fmt.Errorf("each layer must be a table, got %c", value.Type().String())
			return
		}

		layerTable := value.(*lua.LTable)

		var layer *Layer
		layer, err = c.newLayer(layerTable)
		if err != nil {
			return
		}

		layers = append(layers, layer)
	})

	return
}

func (c *Core) getActions(table *lua.LTable) (actions []*Action, err error) {
	field := table.RawGetString(FieldActions)

	if field.Type() == lua.LTNil {
		return nil, nil
	} else if field.Type() != lua.LTTable {
		return nil, fmt.Errorf("actions must be a table, got %c", field.Type().String())
	}

	table = field.(*lua.LTable)

	table.ForEach(func(_, value lua.LValue) {
		if err != nil {
			return
		}

		if value.Type() != lua.LTTable {
			err = fmt.Errorf("each action must be a table, got %c", value.Type().String())
			return
		}

		actionTable := value.(*lua.LTable)

		var action *Action
		action, err = c.newAction(actionTable)
		if err != nil {
			return
		}

		actions = append(actions, action)
	})

	return
}

func New(L *lua.LState, r io.Reader) (*Core, error) {
	lfunc, err := L.Load(r, Module)
	if err != nil {
		return nil, err
	}

	L.Push(lfunc)

	err = L.CallByParam(lua.P{
		Fn:      lfunc,
		NRet:    1,
		Protect: true,
	})

	if err != nil {
		return nil, errCore(err)
	}

	module := L.Get(-1)
	core := &Core{state: L}

	table, ok := module.(*lua.LTable)
	if !ok {
		return nil, errCore(fmt.Errorf("core module must return a table, got %s", module.Type().String()))
	}

	searchTable := table.RawGetString(FieldSearch)
	if searchTable.Type() != lua.LTNil {
		core.search, err = core.newSearch(searchTable.(*lua.LTable))
		if err != nil {
			return nil, errCore(err)
		}
	}

	core.layers, err = core.getLayers(table)
	if err != nil {
		return nil, errCore(err)
	}

	core.actions, err = core.getActions(table)
	if err != nil {
		return nil, errCore(err)
	}

	if !core.HasLayers() && !core.HasSearch() {
		return nil, errCore(fmt.Errorf("core must implement `search` handler or have more than 0 layers"))
	}

	return core, nil
}
