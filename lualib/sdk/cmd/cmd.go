package cmd

import (
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/extensions"
	"github.com/vivi-app/vivi/luautil"
	lua "github.com/yuin/gopher-lua"
	"os/exec"
)

func New(L *lua.LState) *lua.LTable {
	registerCommandType(L)
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"new": newCommand,
	})
}

func newCommand(L *lua.LState) int {
	command := L.CheckString(1)

	programs := L.Context().Value(true).(extensions.ExtensionContainer).Passport().Requirements.Programs

	if !lo.Contains(programs, command) {
		L.RaiseError("command `%s` is not allowed because it is not in the list of allowed programs %s in the extension's passport", command, programs)
		return 0
	}

	cmd := exec.Command(command)
	pushCommand(L, cmd)
	return 1
}

func checkArgs(L *lua.LState, n int) []string {
	args := L.Get(n)

	var (
		argsSlice []string
		err       error
	)

	switch args.Type() {
	case lua.LTString:
		words, err := shellquote.Split(args.String())
		if err != nil {
			L.RaiseError(err.Error())
			return nil
		}

		argsSlice = words
	case lua.LTTable:
		args.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			if err != nil {
				return
			}

			if value.Type() != lua.LTString {
				err = fmt.Errorf("cmd.run: args must be a table of strings")
				return
			}

			argsSlice = append(argsSlice, value.String())
		})

		if err != nil {
			L.RaiseError(err.Error())
		}
	case lua.LTNil:
		// do nothing
	default:
		L.RaiseError("cmd.run: args must be a string or a table of strings")
	}

	return argsSlice
}

func commandSetArgs(L *lua.LState) int {
	cmd := checkCommand(L, 1)
	args := checkArgs(L, 2)

	// do not change the command itself
	cmd.Args = append([]string{cmd.Args[0]}, args...)
	return 0
}

func commandArgs(L *lua.LState) int {
	cmd := checkCommand(L, 1)
	table, err := luautil.ToLValue(L, cmd.Args[1:])
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(table)
	return 1
}

func commandOutput(L *lua.LState) int {
	cmd := checkCommand(L, 1)

	out, err := cmd.Output()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(out))
	return 1
}

func commandRun(L *lua.LState) int {
	cmd := checkCommand(L, 1)

	err := cmd.Run()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}
