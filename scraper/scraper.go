package scraper

import (
	"fmt"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/lualib"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
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

func Parse(r io.Reader) (*Scraper, error) {
	//proto, err := Compile(r)
	//if err != nil {
	//	return nil, err
	//}
	//
	L := lua.NewState()
	lualib.Preload(L)

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

	module := L.Get(-1)
	theScraper := &Scraper{}

	// get script return value
	table, ok := module.(*lua.LTable)
	if !ok {
		fmt.Printf(module.Type().String(), module.String())
		return nil, fmt.Errorf("scraper module must return a table")
	}

	functionSearch := table.RawGet(lua.LString(constant.FunctionSearch))
	if functionSearch.Type() == lua.LTFunction {
		theScraper.functionSearch = functionSearch.(*lua.LFunction)
	}

	theScraper.state = L
	return theScraper, nil
}

func FromPath(path string) (*Scraper, error) {
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

	return Parse(file)
}

func Compile(r io.Reader) (*lua.FunctionProto, error) {
	chunk, err := parse.Parse(r, constant.ModuleScraper)

	if err != nil {
		return nil, err
	}

	proto, err := lua.Compile(chunk, constant.ModuleScraper)
	if err != nil {
		return nil, err
	}

	return proto, nil
}
