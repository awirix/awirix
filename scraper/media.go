package scraper

import (
	"github.com/pkg/errors"
	"github.com/vivi-app/lua"
)

type Media struct {
	scraper  *Scraper
	internal lua.LValue
	*media
}

type media struct {
	Title       string
	Description string
	Info        *lua.LFunction
}

func (i *Media) String() string {
	return i.Title
}

func (i *Media) Description() string {
	return i.media.Description
}

func (i *Media) Value() lua.LValue {
	return i.internal
}

func (i *Media) HasInfo() bool {
	return i.media.Info != nil
}

func (i *Media) Info() (string, error) {
	err := i.scraper.state.CallByParam(lua.P{
		Fn:      i.media.Info,
		NRet:    1,
		Protect: true,
	}, i.internal, i.scraper.context)

	if err != nil {
		return "", err
	}

	info := i.scraper.state.CheckString(-1)
	i.scraper.state.Pop(1)
	return info, nil
}

func (s *Scraper) newMedia(table *lua.LTable) (*Media, error) {
	aux := &media{}
	err := tableMapper.Map(table, aux)
	if err != nil {
		return nil, errors.Wrap(err, "media")
	}

	if aux.Title == "" {
		return nil, errors.Wrap(ErrMissingTitle, "media")
	}

	return &Media{scraper: s, media: aux, internal: table}, nil
}
