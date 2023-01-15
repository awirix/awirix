package scraper

import (
	"github.com/awirix/lua"
	"github.com/pkg/errors"
	"strings"
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

func errMedia(err error) error {
	return errors.Wrap(err, "media")
}

func (m *Media) String() string {
	return m.Title
}

func (m *Media) Description() string {
	return m.media.Description
}

func (m *Media) Value() lua.LValue {
	return m.internal
}

func (m *Media) HasInfo() bool {
	return m.media.Info != nil
}

func (m *Media) Info() (string, error) {
	err := m.scraper.state.CallByParam(lua.P{
		Fn:      m.media.Info,
		NRet:    1,
		Protect: true,
	}, m.internal, m.scraper.context)

	if err != nil {
		return "", errMedia(err)
	}

	info := m.scraper.state.CheckAny(-1)
	if info.Type() != lua.LTString {
		return "", errMedia(errors.New("info must be a string"))
	}

	m.scraper.state.Pop(1)
	return info.String(), nil
}

func (s *Scraper) newMedia(table *lua.LTable) (*Media, error) {
	aux := &media{}
	err := tableMapper.Map(table, aux)

	if err != nil {
		return nil, errMedia(err)
	}

	if aux.Title == "" {
		return nil, errMedia(ErrMissingTitle)
	}

	aux.Title = strings.TrimSpace(aux.Title)
	aux.Description = strings.TrimSpace(aux.Description)

	return &Media{scraper: s, media: aux, internal: table}, nil
}
