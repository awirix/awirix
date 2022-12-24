package scraper

import (
	"fmt"
	lua "github.com/vivi-app/lua"
)

type Media struct {
	internal lua.LValue
	display  string
	about    string
}

func (i *Media) String() string {
	return i.display
}

func (i *Media) Value() lua.LValue {
	return i.internal
}

func (i *Media) About() string {
	return i.about
}

func newMedia(table *lua.LTable) (*Media, error) {
	value := table.RawGetString(FieldTitle)
	display, ok := value.(lua.LString)
	if !ok {
		return nil, fmt.Errorf("invalid media `%s` field: expected 'string' got '%s'", FieldTitle, value.Type().String())
	}

	value = table.RawGetString(FieldDescription)
	about, ok := value.(lua.LString)
	if !ok && value.Type() != lua.LTNil {
		return nil, fmt.Errorf("invalid media `%s` field: expected 'string|nil' got '%s'", FieldDescription, value.Type().String())
	}

	return &Media{
		internal: table,
		display:  string(display),
		about:    string(about),
	}, nil
}
