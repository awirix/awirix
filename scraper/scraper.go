package scraper

import (
	"fmt"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/vm"
	lua "github.com/yuin/gopher-lua"
	"io"
	"path/filepath"
)

type Scraper struct {
	state *lua.LState

	functionSearch *lua.LFunction
}

func (s *Scraper) HasSearch() bool {
	return s.functionSearch != nil
}

func New(path string, r io.Reader) (*Scraper, error) {
	L := vm.New(path)

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

func NewFromPath(path string) (*Scraper, error) {
	isDir, err := filesystem.Api().IsDir(path)
	if err != nil {
		return nil, err
	}

	if isDir {
		path = filepath.Join(path, constant.Scraper)
	}

	file, err := filesystem.Api().Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return New(filepath.Dir(path), file)
}
