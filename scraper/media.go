package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
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
	value := table.RawGetString(FieldDisplay)
	display, ok := value.(lua.LString)
	if !ok {
		return nil, fmt.Errorf("invalid media `%s` field: expected 'string' got '%s'", FieldDisplay, value.Type().String())
	}

	value = table.RawGetString(FieldAbout)
	about, ok := value.(lua.LString)
	if !ok && value.Type() != lua.LTNil {
		return nil, fmt.Errorf("invalid media `%s` field: expected 'string|nil' got '%s'", FieldAbout, value.Type().String())
	}

	return &Media{
		internal: table,
		display:  string(display),
		about:    string(about),
	}, nil
}
