package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

type Intermediate struct {
	internal lua.LValue
	display  string
}

func (i *Intermediate) String() string {
	return i.display
}

func (i *Intermediate) Value() lua.LValue {
	return i.internal
}

func newIntermediate(table *lua.LTable) (*Intermediate, error) {
	value := table.RawGet(lua.LString(FieldDisplay))
	display, ok := value.(lua.LString)
	if !ok {
		return nil, fmt.Errorf("invalid intermediate `%s` field: expected 'display' got '%s'", FieldDisplay, value.Type().String())
	}

	return &Intermediate{
		internal: table,
		display:  string(display),
	}, nil
}
