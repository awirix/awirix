package scraper

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
)

type Scraper struct {
	state *lua.LState

	functionSearch  *lua.LFunction
	functionExplore *lua.LFunction
}

func (s *Scraper) HasSearch() bool {
	return s.functionSearch != nil
}

func (s *Scraper) HasExplore() bool {
	return s.functionExplore != nil
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

	module := L.Get(L.GetTop())
	theScraper := &Scraper{}

	table, ok := module.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("scraper module must return a table, got %s", module.Type().String())
	}

	functionSearch := table.RawGet(lua.LString(FunctionSearch))
	if functionSearch.Type() == lua.LTFunction {
		theScraper.functionSearch = functionSearch.(*lua.LFunction)
	}

	functionExplore := table.RawGet(lua.LString(FunctionExplore))
	if functionExplore.Type() == lua.LTFunction {
		theScraper.functionExplore = functionExplore.(*lua.LFunction)
	}

	if !theScraper.HasExplore() && !theScraper.HasSearch() {
		return nil, fmt.Errorf("scraper module must return at least one of the functions `%s` or `%s`", FunctionSearch, FunctionExplore)
	}

	theScraper.state = L
	return theScraper, nil
}
