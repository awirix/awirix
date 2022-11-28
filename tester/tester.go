package tester

import (
	"fmt"
	"github.com/vivi-app/vivi/constant"
	lua "github.com/yuin/gopher-lua"
	"io"
)

type Tester struct {
	state        *lua.LState
	functionTest *lua.LFunction
}

func New(L *lua.LState, r io.Reader) (*Tester, error) {
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
