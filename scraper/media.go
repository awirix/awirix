package scraper

import (
	"github.com/pkg/errors"
	"github.com/vivi-app/lua"
)

type Media struct {
	internal    lua.LValue
	Title       string
	Description string
}

func (i *Media) String() string {
	return i.Title
}

func (i *Media) Value() lua.LValue {
	return i.internal
}

func newMedia(table *lua.LTable) (*Media, error) {
	media := &Media{internal: table}
	err := tableMapper.Map(table, media)

	if err != nil {
		return nil, errors.Wrap(err, "media")
	}

	if media.Title == "" {
		return nil, errors.Wrap(ErrMissingTitle, "media")
	}

	return media, nil
}
