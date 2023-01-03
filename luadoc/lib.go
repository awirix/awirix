package luadoc

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
	"strings"
)

type Lib struct {
	Name        string
	Description string
	Vars        []*Var
	Funcs       []*Func
	Classes     []*Class
	Libs        []*Lib
}

func (l *Lib) Value(state *lua.LState) *lua.LTable {
	var (
		vars  = make(map[string]lua.LValue)
		funcs = make(map[string]lua.LGFunction)
	)

	for _, f := range l.Funcs {
		if f.Value == nil {
			panic(fmt.Sprintf("function %s has no value", f.Name))
		}

		funcs[f.Name] = f.Value
	}

	for _, v := range l.Vars {
		if v.Value == nil {
			panic(fmt.Sprintf("variable %s has no value", v.Name))
		}

		vars[v.Name] = v.Value
	}

	for _, l := range l.Libs {
		vars[l.Name] = l.Value(state)
	}

	for _, c := range l.Classes {
		var methods = make(map[string]lua.LGFunction)
		for _, m := range c.Methods {
			if m.Value == nil {
				panic(fmt.Sprintf("method %s of class %s has no value", m.Name, c.Name))
			}
			methods[m.Name] = m.Value
		}

		mt := state.NewTypeMetatable(c.Name)
		state.SetField(mt, "__index", state.SetFuncs(state.NewTable(), methods))
	}

	return luautil.NewTable(state, vars, funcs)
}

func (l *Lib) LuaDoc() string {
	var b strings.Builder
	lo.Must0(templateLuaDocLib.Execute(&b, l))

	return b.String()
}

func (l *Lib) Loader() func(state *lua.LState) int {
	return func(state *lua.LState) int {
		libs := state.NewTable()

		for _, lib := range l.Libs {
			state.SetField(libs, lib.Name, lib.Value(state))
		}

		state.Push(libs)
		return 1
	}
}
