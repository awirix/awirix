package scraper

import (
	"fmt"
	"github.com/vivi-app/lua"
)

func errOneOfRequired(functions ...*lua.LFunction) error {
	// TODO: add names of the handler
	return fmt.Errorf("at least one of the following functions is required: %v", functions)
}

func errNotAFunction(name string, val lua.LValue) error {
	return fmt.Errorf("scraper module must return a handler `%s`, got %s", name, val.Type().String())
}

func getFunctionFromTable(table *lua.LTable, name string, required bool) (*lua.LFunction, error) {
	function := table.RawGetString(name)

	if function.Type() == lua.LTFunction {
		return function.(*lua.LFunction), nil
	} else if function.Type() == lua.LTNil && !required {
		return nil, nil
	}

	return nil, errNotAFunction(name, function)
}
