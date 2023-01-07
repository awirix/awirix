package cmd

import (
	"fmt"
	"github.com/awirix/awirix/extensions"
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/awirix/luautil"
	"github.com/awirix/lua"
	"github.com/kballard/go-shellquote"
	"github.com/samber/lo"
	"os/exec"
)

func Lib() *luadoc.Lib {
	classCmd := &luadoc.Class{
		Name:        "command",
		Description: "Command object that is used to execute external programs.",
		Methods: []*luadoc.Method{
			{
				Name:        "run",
				Description: "Runs the command and returns an error if it fails.",
				Value:       commandRun,
				Params:      []*luadoc.Param{},
				Returns: []*luadoc.Param{
					{
						Name:        "error",
						Description: "Error if the command fails.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "output",
				Description: "Runs the command and returns the output and an error if it fails.",
				Value:       commandOutput,
				Returns: []*luadoc.Param{
					{
						Name:        "output",
						Description: "Output of the command.",
						Type:        luadoc.String,
					},
					{
						Name:        "error",
						Description: "Error if the command fails.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "get_args",
				Description: "Returns the arguments of the command.",
				Value:       commandGetArgs,
				Params:      []*luadoc.Param{},
				Returns: []*luadoc.Param{
					{
						Name:        "args",
						Description: "Arguments of the command.",
						Type:        luadoc.List(luadoc.String),
					},
				},
			},
			{
				Name:        "set_args",
				Description: "Sets the arguments of the command.",
				Value:       commandArgs,
				Params: []*luadoc.Param{
					{
						Name:        "args",
						Description: "Arguments of the command.",
						Type:        luadoc.List(luadoc.String),
					},
				},
			},
		},
	}

	return &luadoc.Lib{
		Name:        "cmd",
		Description: "The `cmd` library provides a way to run external programs.",
		Funcs: []*luadoc.Func{
			{
				Name:        "new",
				Description: "Creates a new command object. The command object is a wrapper around the `exec.Cmd` object from the Go standard library.",
				Value:       newCommand,
				Params: []*luadoc.Param{
					{
						Name:        "command",
						Description: "The command to run. This must be a string and must be in the list of allowed programs in the extension's passport.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "command",
						Description: "The command object.",
						Type:        classCmd.Name,
					},
				},
			},
		},
		Classes: []*luadoc.Class{classCmd},
	}
}

func newCommand(L *lua.LState) int {
	command := L.CheckString(1)

	programs := L.Context().Value("extension").(extensions.ExtensionContainer).Passport().Requirements.Programs

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

func commandArgs(L *lua.LState) int {
	cmd := checkCommand(L, 1)
	args := checkArgs(L, 2)

	// do not change the command itself
	cmd.Args = append([]string{cmd.Args[0]}, args...)
	return 0
}

func commandGetArgs(L *lua.LState) int {
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

const commandTypeName = "command"

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
