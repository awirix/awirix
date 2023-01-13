package levenshtein

import (
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/lua"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "levenshtein",
		Description: "Levenshtein distance algorithm",
		Funcs: []*luadoc.Func{
			{
				Name:        "distance",
				Description: "Compute Levenshtein distance between two strings",
				Value:       distance,
				Params: []*luadoc.Param{
					{
						Name: "s1",
						Type: luadoc.String,
					},
					{
						Name: "s2",
						Type: luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "distance",
						Description: "Levenshtein distance between s1 and s2",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "closest",
				Description: "Find string with highest similarity from strings list",
				Value:       closest,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "String to compare with",
						Type:        luadoc.String,
					},
					{
						Name:        "strings",
						Description: "List of strings to compare with",
						Type:        luadoc.List(luadoc.String),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "closest",
						Description: "String with highest similarity",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func distance(L *lua.LState) int {
	s1 := L.CheckString(1)
	s2 := L.CheckString(2)

	L.Push(lua.LNumber(levenshtein.Distance(s1, s2)))

	return 1
}

func closest(L *lua.LState) int {
	s := L.CheckString(1)
	strings := L.CheckTable(2)

	var (
		closest  string
		distance int
	)

	strings.ForEach(func(key, value lua.LValue) {
		s2 := value.String()

		d := levenshtein.Distance(s, s2)

		if closest == "" || d < distance {
			closest = s2
			distance = d
		}
	})

	L.Push(lua.LString(closest))

	return 1
}
