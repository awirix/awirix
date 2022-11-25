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
	HasSearch bool
	state     *lua.LState

	functionSearch *lua.LFunction
}

func Parse(r io.Reader) (*Scraper, error) {
	proto, err := Compile(r)
	if err != nil {
		return nil, err
	}

	L := lua.NewState()
	lualib.Preload(L)

	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	err = L.PCall(0, 1, nil)
	if err != nil {
		return nil, err
	}

	theScraper := &Scraper{}

	// get script return value
	module := L.GetGlobal(constant.ModuleScraper)
	table, ok := module.(*lua.LTable)
	if !ok {
		fmt.Printf(module.Type().String(), module.String())
		return nil, fmt.Errorf("invalid scraper")
	}

	hasSearch := table.RawGet(lua.LString(constant.FieldHasSearch))
	if hasSearch.Type() == lua.LTBool {
		theScraper.HasSearch = bool(hasSearch.(lua.LBool))
	}

	functionSearch := table.RawGet(lua.LString(constant.FunctionSearch))
	if functionSearch.Type() == lua.LTFunction {
		theScraper.HasSearch = true
		theScraper.functionSearch = functionSearch.(*lua.LFunction)
	} else {
		theScraper.HasSearch = false
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
