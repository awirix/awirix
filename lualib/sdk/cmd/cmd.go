package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
	"os/exec"
)

func New(L *lua.LState) *lua.LTable {
	return util.NewTable(L, nil, map[string]lua.LGFunction{
		"run": run,
	})
}

func run(L *lua.LState) int {
	cmd := L.CheckString(1)
	args := L.CheckTable(2)

	var (
		argsSlice []string
		err       error
	)
	args.ForEach(func(key, value lua.LValue) {
		if value.Type() != lua.LTString {
			err = fmt.Errorf("cmd.run: args must be a table of strings")
			return
		}

		argsSlice = append(argsSlice, value.String())
	})

	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	if viper.GetBool(constant.ExtensionsSafeMode) {
		L.Push(lua.LNil)
		L.Push(lua.LString("command execution is disabled in safe mode"))
		return 2
	}

	out, err := exec.Command(cmd, argsSlice...).Output()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(string(out)))
	return 1
}
