package tester

import (
	"fmt"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/vm"
	lua "github.com/yuin/gopher-lua"
	"io"
	"path/filepath"
)

type Tester struct {
	state        *lua.LState
	functionTest *lua.LFunction
}

func New(path string, r io.Reader) (*Tester, error) {
	L := vm.New(path)

	lfunc, err := L.Load(r, constant.ModuleTest)
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
	theTester := &Tester{}

	// get script return value
	table, ok := module.(*lua.LTable)
	if !ok {
		fmt.Printf(module.Type().String(), module.String())
		return nil, fmt.Errorf("tester module must return a table")
	}

	functionTest := table.RawGet(lua.LString(constant.FunctionTest))
	if functionTest.Type() == lua.LTFunction {
		theTester.functionTest = functionTest.(*lua.LFunction)
	} else {
		return nil, fmt.Errorf("tester module must have a %s function", constant.FunctionTest)
	}

	theTester.state = L
	return theTester, nil
}

func NewFromPath(path string) (*Tester, error) {
	isDir, err := filesystem.Api().IsDir(path)
	if err != nil {
		return nil, err
	}

	if isDir {
		path = filepath.Join(path, constant.Tester)
	}

	file, err := filesystem.Api().Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return New(filepath.Dir(path), file)
}
