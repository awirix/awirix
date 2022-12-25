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

func getNoun(table *lua.LTable) (*Noun, error) {
	noun := &Noun{}

	value := table.RawGetString(FieldNoun)
	if value.Type() == lua.LTNil {
		noun.singular = "media"
		return noun, nil
	}

	if value.Type() != lua.LTTable {
		return nil, fmt.Errorf("noun: must be a table, got %s", value.Type().String())
	}

	value = table.RawGetString(FieldSingular)

	if value.Type() != lua.LTNil {
		if value.Type() != lua.LTString {
			return nil, fmt.Errorf("noun: singular must be a string, got %s", value.Type().String())
		}

		noun.singular = string(value.(lua.LString))
	}

	value = table.RawGetString(FieldPlural)

	if value.Type() != lua.LTNil {
		if value.Type() != lua.LTString {
			return nil, fmt.Errorf("noun: plural must be a string, got %s", value.Type().String())
		}

		noun.plural = string(value.(lua.LString))
	}

	return noun, nil
}

func (s *Scraper) newLayer(table *lua.LTable) (*Layer, error) {
	value := table.RawGetString(FieldTitle)

	if value.Type() != lua.LTString {
		return nil, fmt.Errorf("layer must have a string title, got %s", value.Type().String())
	}

	title := string(value.(lua.LString))

	value = table.RawGetString(FieldHandler)

	if value.Type() != lua.LTFunction {
		return nil, fmt.Errorf("layer must have a handler handler, got %s", value.Type().String())
	}

	handler := value.(*lua.LFunction)

	noun, err := getNoun(table)
	if err != nil {
		return nil, err
	}

	return &Layer{
		title:   title,
		handler: handler,
		scraper: s,
		noun:    noun,
	}, nil
}

func (s *Scraper) getSearch(table *lua.LTable) (*Search, error) {
	value := table.RawGetString(FieldSearch)

	if value.Type() == lua.LTNil {
		return nil, nil
	}

	if value.Type() != lua.LTTable {
		return nil, fmt.Errorf("search must be a table, got %s", value.Type().String())
	}

	table = value.(*lua.LTable)

	value = table.RawGetString(FieldTitle)
	if value.Type() != lua.LTString {
		return nil, fmt.Errorf("search must have a string title, got %s", value.Type().String())
	}

	title := string(value.(lua.LString))

	value = table.RawGetString(FieldHandler)
	if value.Type() != lua.LTFunction {
		return nil, fmt.Errorf("search must have a handler function, got %s", value.Type().String())
	}

	handler := value.(*lua.LFunction)

	noun, err := getNoun(table)
	if err != nil {
		return nil, err
	}

	var subtitle string
	value = table.RawGetString(FieldSubtitle)
	if value.Type() != lua.LTNil {
		if value.Type() != lua.LTString {
			return nil, fmt.Errorf("search: subtitle must be a string, got %s", value.Type().String())
		}

		subtitle = string(value.(lua.LString))
	}

	var placeholder string
	value = table.RawGetString(FieldPlaceholder)
	if value.Type() != lua.LTNil {
		if value.Type() != lua.LTString {
			return nil, fmt.Errorf("search: placeholder must be a string, got %s", value.Type().String())
		}

		placeholder = string(value.(lua.LString))
	}

	return &Search{
		title:       title,
		handler:     handler,
		scraper:     s,
		noun:        noun,
		subtitle:    subtitle,
		placeholder: placeholder,
	}, nil
}
