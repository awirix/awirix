package tester

import lua "github.com/yuin/gopher-lua"

func (t *Tester) Test() error {
	err := t.state.CallByParam(lua.P{
		Fn:      t.functionTest,
		NRet:    0,
		Protect: true,
	})

	if err != nil {
		return err
	}

	return nil
}