package fuzzy

import (
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/lua"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "fuzzy",
		Description: "Fuzzy search algorithm",
		Funcs: []*luadoc.Func{
			{
				Name:        "find",
				Description: "Find string with highest similarity from strings list",
				Value:       find,
				Params: []*luadoc.Param{
					{
						Name:        "pattern",
						Description: "Pattern to compare with",
						Type:        luadoc.String,
					},
					{
						Name:        "data",
						Description: "Data to search in",
						Type:        luadoc.List(luadoc.String),
					},
				},
			},
		},
	}
}

func find(L *lua.LState) int {
	pattern := L.CheckString(1)
	strings := L.CheckTable(2)

	var list []string

	strings.ForEach(func(key lua.LValue, value lua.LValue) {
		if value.Type() != lua.LTString {
			L.RaiseError("fuzzy.find: table must contain only strings")
			return
		}

		list = append(list, value.String())
	})

	found := fuzzy.Find(pattern, list)

	ret := L.NewTable()
	for _, f := range found {
		ret.Append(lua.LString(f))
	}

	L.Push(ret)
	return 1
}
