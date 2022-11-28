package scraper

import (
	"fmt"
	"github.com/vivi-app/vivi/constant"
	lua "github.com/yuin/gopher-lua"
	"io"
)

type Scraper struct {
	state *lua.LState

	functionSearch *lua.LFunction
}

func (s *Scraper) HasSearch() bool {
	return s.functionSearch != nil
}

func New(L *lua.LState, r io.Reader) (*Scraper, error) {
	lfunc, err := L.Load(r, constant.ModuleScraper)
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

	functionSearch := table.RawGet(lua.LString(constant.FunctionSearch))
	if functionSearch.Type() == lua.LTFunction {
		theScraper.functionSearch = functionSearch.(*lua.LFunction)
	}

	theScraper.state = L
	return theScraper, nil
}
