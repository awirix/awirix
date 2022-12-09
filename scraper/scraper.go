package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
)

type Scraper struct {
	state *lua.LState

	functionSearch   *lua.LFunction
	functionExplore  *lua.LFunction
	functionPrepare  *lua.LFunction
	functionStream   *lua.LFunction
	functionDownload *lua.LFunction
	progress         *lua.LFunction
}

func (s *Scraper) HasSearch() bool {
	return s.functionSearch != nil
}

func (s *Scraper) HasExplore() bool {
	return s.functionExplore != nil
}

func (s *Scraper) HasStream() bool {
	return s.functionStream != nil
}

func (s *Scraper) HasDownload() bool {
	return s.functionDownload != nil
}

func (s *Scraper) SetProgress(progress func(string)) {
	s.progress = s.state.NewFunction(func(L *lua.LState) int {
		progress(L.ToString(1))
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
	function := table.RawGet(lua.LString(name))

	if function.Type() == lua.LTFunction {
		return function.(*lua.LFunction), nil
	} else if function.Type() == lua.LTNil && !required {
		return nil, nil
	}

	return nil, errNotAFunction(name, function)
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

	theScraper.functionExplore, err = getFunctionFromTable(table, FunctionExplore, false)
	if err != nil {
		return nil, err
	}

	if !theScraper.HasExplore() && !theScraper.HasSearch() {
		return nil, errOneOfRequired(theScraper.functionSearch, theScraper.functionExplore)
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
