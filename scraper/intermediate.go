package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

type Intermediate struct {
	internal lua.LValue
	string   string
}

func (i *Intermediate) String() string {
	return i.string
}

func (i *Intermediate) Value() lua.LValue {
	return i.internal
}

func newIntermediate(table *lua.LTable) (*Intermediate, error) {
	value := table.RawGet(lua.LString("string"))
	str, ok := value.(lua.LString)
	if !ok {
		return nil, fmt.Errorf("invalid intermediate 'string' field: expected 'string' got '%s'", value.Type().String())
	}

	return &Intermediate{
		internal: table,
		string:   string(str),
	}, nil
}
