package core

import (
	"fmt"
	"github.com/awirix/lua"
)

func (c *Core) checkMedia() (*Media, error) {
	ret := c.state.Get(-1)
	c.state.Pop(1)

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid return value: expected 'table' got '%c'", ret.Type().String())
	}

	media, err := c.newMedia(table)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (c *Core) checkMediaSlice() ([]*Media, error) {
	ret := c.state.Get(-1)
	c.state.Pop(1)

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid return value: expected 'table' got '%c'", ret.Type().String())
	}

	var (
		items  = make([]lua.LValue, 0)
		medias = make([]*Media, 0)
	)

	table.ForEach(func(_, value lua.LValue) {
		items = append(items, value)
	})

	for _, item := range items {
		table, ok := item.(*lua.LTable)
		if !ok {
			return nil, fmt.Errorf("invalid value in returned table: expected 'table' got '%c'", item.Type().String())
		}

		media, err := c.newMedia(table)
		if err != nil {
			return nil, err
		}

		medias = append(medias, media)
	}

	return medias, nil
}

func (c *Core) Search() *Search {
	return c.search
}

func (c *Core) Layers() []*Layer {
	return c.layers
}

func (c *Core) Actions() []*Action {
	return c.actions
}
