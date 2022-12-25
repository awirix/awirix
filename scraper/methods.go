package scraper

import (
	"fmt"
	"github.com/vivi-app/lua"
)

func (s *Scraper) checkMedia() (*Media, error) {
	ret := s.state.Get(-1)
	s.state.Pop(1)

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid return value: expected 'table' got '%s'", ret.Type().String())
	}

	media, err := newMedia(table)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s *Scraper) checkMediaSlice() ([]*Media, error) {
	ret := s.state.Get(-1)
	s.state.Pop(1)

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("invalid return value: expected 'table' got '%s'", ret.Type().String())
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
			return nil, fmt.Errorf("invalid value in returned table: expected 'table' got '%s'", item.Type().String())
		}

		media, err := newMedia(table)
		if err != nil {
			return nil, err
		}

		medias = append(medias, media)
	}

	return medias, nil
}

func (s *Scraper) Search() *Search {
	if !s.HasSearch() {
		panic("scraper does not have a search handler")
	}

	return s.search
}

func (s *Scraper) Layers() []*Layer {
	return s.layers
}

func (s *Scraper) Actions() []*Action {
	return s.actions
}
