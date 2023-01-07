package tester

import "github.com/awirix/lua"

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
