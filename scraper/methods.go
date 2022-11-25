package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func (s *Scraper) Search(query string) ([]*Intermediate, error) {
	err := s.state.CallByParam(lua.P{
		Fn:      s.functionSearch,
		NRet:    1,
		Protect: true,
	}, lua.LString(query))

	if err != nil {
		return nil, err
	}

	ret := s.state.Get(-1)
	s.state.Pop(1)

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid search return value: expected 'table' got '%s'", ret.Type().String())
	}

	var (
		items         = make([]lua.LValue, 0)
		intermediates = make([]*Intermediate, 0)
	)

	table.ForEach(func(_, value lua.LValue) {
		items = append(items, value)
	})

	for _, item := range items {
		table, ok := item.(*lua.LTable)
		if !ok {
			return nil, fmt.Errorf("invalid search return value: expected 'table' got '%s'", item.Type().String())
		}

		intermediate, err := newIntermediate(table)
		if err != nil {
			return nil, err
		}

		intermediates = append(intermediates, intermediate)
	}

	return intermediates, nil
}
