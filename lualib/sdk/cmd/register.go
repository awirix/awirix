package cmd

import (
	lua "github.com/vivi-app/lua"
	"os/exec"
)

const commandTypeName = "command"

func registerCommandType(L *lua.LState) {
	mt := L.NewTypeMetatable(commandTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), commandMethods))
}

var commandMethods = map[string]lua.LGFunction{
	"run":      commandRun,
	"output":   commandOutput,
	"set_args": commandArgs,
	"args":     commandGetArgs,
}

func checkCommand(L *lua.LState, n int) *exec.Cmd {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*exec.Cmd); ok {
		return v
	}
	L.ArgError(1, "command expected")
	return nil
}

func pushCommand(L *lua.LState, cmd *exec.Cmd) {
	ud := L.NewUserData()
	ud.Value = cmd
	L.SetMetatable(ud, L.GetTypeMetatable(commandTypeName))
	L.Push(ud)
}
