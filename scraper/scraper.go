package scraper

import (
	"fmt"
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/log"
	"io"
)

type Scraper struct {
	state *lua.LState

	search  *Search
	layers  []*Layer
	actions []*Action

	progress *lua.LFunction
}

func (s *Scraper) HasSearch() bool {
	return s.search != nil
}

func (s *Scraper) HasLayers() bool {
	return s.layers != nil || len(s.layers) > 0
}

func (s *Scraper) HasActions() bool {
	return s.actions != nil || len(s.actions) > 0
}

func (s *Scraper) SetProgress(progress func(string)) {
	s.progress = s.state.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		progress(msg)
		log.Tracef("progress: %s", msg)
		return 0
	})
}

func (s *Scraper) getLayers(table *lua.LTable) (layers []*Layer, err error) {
	field := table.RawGetString(FieldLayers)

	if field.Type() == lua.LTNil {
		return nil, nil
	} else if field.Type() != lua.LTTable {
		return nil, fmt.Errorf("layers must be a table, got %s", field.Type().String())
	}

	table = field.(*lua.LTable)

	table.ForEach(func(_, value lua.LValue) {
		if err != nil {
			return
		}

		if value.Type() != lua.LTTable {
			err = fmt.Errorf("each layer must be a table, got %s", value.Type().String())
			return
		}

		layerTable := value.(*lua.LTable)

		var layer *Layer
		layer, err = s.newLayer(layerTable)
		if err != nil {
			return
		}

		layers = append(layers, layer)
	})

	return
}

func (s *Scraper) getActions(table *lua.LTable) (actions []*Action, err error) {
	field := table.RawGetString(FieldActions)

	if field.Type() == lua.LTNil {
		return nil, nil
	} else if field.Type() != lua.LTTable {
		return nil, fmt.Errorf("actions must be a table, got %s", field.Type().String())
	}

	table = field.(*lua.LTable)

	table.ForEach(func(_, value lua.LValue) {
		if err != nil {
			return
		}

		if value.Type() != lua.LTTable {
			err = fmt.Errorf("each action must be a table, got %s", value.Type().String())
			return
		}

		actionTable := value.(*lua.LTable)

		var action *Action
		action, err = s.newAction(actionTable)
		if err != nil {
			return
		}

		actions = append(actions, action)
	})

	return
}

func New(L *lua.LState, r io.Reader) (*Scraper, error) {
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
		return nil, err
	}

	module := L.Get(-1)
	theScraper := &Scraper{}

	table, ok := module.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("scraper module must return a table, got %s", module.Type().String())
	}

	searchTable := table.RawGetString(FieldSearch)
	if searchTable.Type() != lua.LTNil {
		theScraper.search, err = theScraper.newSearch(searchTable.(*lua.LTable))
		if err != nil {
			return nil, err
		}
	}

	theScraper.layers, err = theScraper.getLayers(table)
	if err != nil {
		return nil, err
	}

	theScraper.actions, err = theScraper.getActions(table)
	if err != nil {
		return nil, err
	}

	if !theScraper.HasLayers() && !theScraper.HasSearch() {
		return nil, fmt.Errorf("scraper must implement `search` handler or have more than 0 layers")
	}

	theScraper.state = L
	return theScraper, nil
}
