package fmt

import (
	"fmt"
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/awirix/luautil"
	"github.com/awirix/lua"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "fmt",
		Description: "Formatting library. Contains various utilities for formatting strings.",
		Funcs: []*luadoc.Func{
			{
				Name:        "sprintf",
				Description: "Formats a string using the given pattern and arguments.",
				Value:       fmtSprintf,
				Params: []*luadoc.Param{
					{
						Name:        "pattern",
						Description: "The pattern to use for formatting.",
						Type:        luadoc.String,
					},
					{
						Name:        "...",
						Description: "The arguments to use for formatting.",
						Type:        luadoc.Any,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "formatted",
						Description: "The formatted string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "sprintln",
				Description: "Formats using the default formats for its operands and returns the resulting string. Spaces are always added between operands and a newline is appended.",
				Value:       fmtSprintln,
				Params: []*luadoc.Param{
					{
						Name:        "...",
						Description: "The arguments to use for formatting.",
						Type:        luadoc.Any,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "formatted",
						Description: "The formatted string.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func fmtSprintf(L *lua.LState) int {
	pattern := L.CheckString(1)
	var args []any

	for i := 2; i <= L.GetTop(); i++ {
		arg, err := luautil.FromLValue(L.Get(i))

		if err != nil {
			L.RaiseError(err.Error())
		}

		args = append(args, arg)
	}

	L.Push(lua.LString(fmt.Sprintf(pattern, args...)))
	return 1
}

func fmtSprintln(L *lua.LState) int {
	var args []any

	for i := 1; i <= L.GetTop(); i++ {
		arg, err := luautil.FromLValue(L.Get(i))

		if err != nil {
			L.RaiseError(err.Error())
		}

		args = append(args, arg)
	}

	L.Push(lua.LString(fmt.Sprintln(args...)))
	return 1
}
