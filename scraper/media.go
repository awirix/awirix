package scraper

import (
	"fmt"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/log"
)

type Media struct {
	internal    lua.LValue
	title       string
	description string
}

func (i *Media) String() string {
	return i.title
}

func (i *Media) Value() lua.LValue {
	return i.internal
}

func (i *Media) Description() string {
	return i.description
}

func newMedia(table *lua.LTable) (*Media, error) {
	var media = &Media{internal: table}

	value := table.RawGetString(FieldTitle)

	if value.Type() != lua.LTNil {
		if value.Type() != lua.LTString {
			return nil, fmt.Errorf("title: must be a string, got %s", value.Type().String())
		}

		media.title = string(value.(lua.LString))
	}

	value = table.RawGetString(FieldDescription)

	if value.Type() != lua.LTNil {
		if value.Type() != lua.LTString {
			return nil, fmt.Errorf("description: must be a string, got %s", value.Type().String())
		}

		log.Info(value.String())
		media.description = string(value.(lua.LString))
	}

	return media, nil
}
