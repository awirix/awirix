package scraper

import (
	"fmt"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/log"
	"io"
)

type Scraper struct {
	state *lua.LState

	functionSearch   *lua.LFunction
	functionPrepare  *lua.LFunction
	functionStream   *lua.LFunction
	functionDownload *lua.LFunction

	layers   []*Layer
	progress *lua.LFunction
}

func (s *Scraper) HasSearch() bool {
	return s.functionSearch != nil
}

func (s *Scraper) HasLayers() bool {
	return s.layers != nil || len(s.layers) > 0
}

func (s *Scraper) HasStream() bool {
	return s.functionStream != nil
}

func (s *Scraper) HasDownload() bool {
	return s.functionDownload != nil
}

func (s *Scraper) SetProgress(progress func(string)) {
	s.progress = s.state.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		progress(msg)
		log.Tracef("progress: %s", msg)
		return 0
	})
}

func errOneOfRequired(functions ...*lua.LFunction) error {
	// TODO: add names of the function
	return fmt.Errorf("at least one of the following functions is required: %v", functions)
}

func errNotAFunction(name string, val lua.LValue) error {
	return fmt.Errorf("scraper module must return a function `%s`, got %s", name, val.Type().String())
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

func getLayers(table *lua.LTable) (layers []*Layer, err error) {
	field := table.RawGetString(FieldLayers)

	if field.Type() == lua.LTNil {
		return nil, nil
	} else if field.Type() != lua.LTTable {
		return nil, fmt.Errorf("layers field must be a table, got %s", field.Type().String())
	}

	table = field.(*lua.LTable)

	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if err != nil {
			return
		}

		if key.Type() != lua.LTString {
			err = fmt.Errorf("each layer name must be a string, got %s", key.Type().String())
			return
		}

		if value.Type() != lua.LTFunction {
			err = fmt.Errorf("each layer must be a function, got %s", value.Type().String())
			return
		}

		layers = append(layers, &Layer{
			Name:        string(key.(lua.LString)),
			luaFunction: value.(*lua.LFunction),
		})
	})

	return
}

func New(L *lua.LState, r io.Reader) (*Scraper, error) {
	lfunc, err := L.Load(r, Module)
	if err != nil {
		return nil, err
	}

	L.Push(lfunc)

	err = L.CallByParam(lua.P{
		Fn:      lfunc,
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return nil, err
	}

	module := L.Get(-1)
	theScraper := &Scraper{}

	table, ok := module.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("scraper module must return a table, got %s", module.Type().String())
	}

	theScraper.functionSearch, err = getFunctionFromTable(table, FunctionSearch, false)
	if err != nil {
		return nil, err
	}

	theScraper.layers, err = getLayers(table)
	if err != nil {
		return nil, err
	}

	if !theScraper.HasLayers() && !theScraper.HasSearch() {
		return nil, fmt.Errorf("scraper must implement `search` function or have more than 0 layers")
	}

	theScraper.functionPrepare, err = getFunctionFromTable(table, FunctionPrepare, true)
	if err != nil {
		return nil, err
	}

	theScraper.functionStream, err = getFunctionFromTable(table, FunctionStream, false)
	if err != nil {
		return nil, err
	}

	theScraper.functionDownload, err = getFunctionFromTable(table, FunctionDownload, false)
	if err != nil {
		return nil, err
	}

	if !theScraper.HasDownload() && !theScraper.HasStream() {
		return nil, errOneOfRequired(theScraper.functionDownload, theScraper.functionStream)
	}

	theScraper.state = L
	return theScraper, nil
}
