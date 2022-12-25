package scraper

import (
	"github.com/vivi-app/gluamapper"
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
	media := &Media{}
	err := gluamapper.Map(table, media)
	if err != nil {
		return nil, err
	}

	media.internal = table
	return media, nil
}
