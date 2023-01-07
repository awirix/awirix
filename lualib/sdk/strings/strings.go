package strings

import (
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/awirix/luautil"
	lua "github.com/awirix/lua"
	"strings"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "strings",
		Description: "Simple functions to manipulate UTF-8 encoded strings.",
		Funcs: []*luadoc.Func{
			{
				Name:        "replace_all",
				Description: "Replaces all instances of old with new in s.",
				Value:       replaceAll,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to replace in.",
						Type:        luadoc.String,
					},
					{
						Name:        "old",
						Description: "The string to replace.",
						Type:        luadoc.String,
					},
					{
						Name:        "new",
						Description: "The string to replace with.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "replace",
				Description: "Replaces the first n instances of old with new in s.",
				Value:       replace,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to replace in.",
						Type:        luadoc.String,
					},
					{
						Name:        "old",
						Description: "The string to replace.",
						Type:        luadoc.String,
					},
					{
						Name:        "new",
						Description: "The string to replace with.",
						Type:        luadoc.String,
					},
					{
						Name:        "n",
						Description: "The number of instances to replace.",
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "split",
				Description: "Splits s around each instance of sep and returns a table of the substrings between those instances.",
				Value:       split,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to split.",
						Type:        luadoc.String,
					},
					{
						Name:        "sep",
						Description: "The string to split around.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "table",
						Description: "A table of the substrings between each instance of sep.",
						Type:        luadoc.Table,
					},
				},
			},
			{
				Name:        "to_lower",
				Description: "Returns a copy of the string s with all Unicode letters mapped to their lower case.",
				Value:       toLower,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to convert.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "to_upper",
				Description: "Returns a copy of the string s with all Unicode letters mapped to their upper case.",
				Value:       toUpper,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to convert.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "trim",
				Description: "Returns a copy of the string s with all leading and trailing Unicode code points contained in cutset removed.",
				Value:       trim,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to trim.",
						Type:        luadoc.String,
					},
					{
						Name:        "cutset",
						Description: "The set of Unicode code points to remove.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "trim_left",
				Description: "Returns a copy of the string s with all leading Unicode code points contained in cutset removed.",
				Value:       trimLeft,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to trim.",
						Type:        luadoc.String,
					},
					{
						Name:        "cutset",
						Description: "The set of Unicode code points to remove.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "trim_right",
				Description: "Returns a copy of the string s with all trailing Unicode code points contained in cutset removed.",
				Value:       trimRight,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to trim.",
						Type:        luadoc.String,
					},
					{
						Name:        "cutset",
						Description: "The set of Unicode code points to remove.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "trim_space",
				Description: "Returns a copy of the string s with all leading and trailing white space removed, as defined by Unicode.",
				Value:       trimSpace,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to trim.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "trim_prefix",
				Description: "Returns a copy of the string s without the provided leading prefix string. If s doesn't start with prefix, s is returned unchanged.",
				Value:       trimPrefix,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to trim.",
						Type:        luadoc.String,
					},
					{
						Name:        "prefix",
						Description: "The prefix to remove.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "trim_suffix",
				Description: "Returns a copy of the string s without the provided trailing suffix string. If s doesn't end with suffix, s is returned unchanged.",
				Value:       trimSuffix,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to trim.",
						Type:        luadoc.String,
					},
					{
						Name:        "suffix",
						Description: "The suffix to remove.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "new",
						Description: "The new string.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "has_prefix",
				Description: "Returns true if the string s begins with prefix.",
				Value:       hasPrefix,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "prefix",
						Description: "The prefix to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "True if the string has the prefix.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "has_suffix",
				Description: "Returns true if the string s ends with suffix.",
				Value:       hasSuffix,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "suffix",
						Description: "The suffix to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "True if the string has the suffix.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "contains",
				Description: "Returns true if the string s contains substr.",
				Value:       contains,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "substr",
						Description: "The substring to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "True if the string contains the substring.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name: "contains_any",
				Description: `Returns true if the string s contains any of the runes in chars.
If chars is the empty string, contains_any returns false.`,
				Value: containsAny,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "chars",
						Description: "The characters to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "True if the string contains any of the characters.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "count",
				Description: "Returns the number of non-overlapping instances of substr in s.",
				Value:       count,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "substr",
						Description: "The substring to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "n",
						Description: "The number of non-overlapping instances of substr in s.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "count_any",
				Description: "Returns the number of non-overlapping instances of any of the runes in chars in s.",
				Value:       countAny,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "chars",
						Description: "The characters to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "n",
						Description: "The number of non-overlapping instances of any of the runes in chars in s.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "equal_fold",
				Description: "Returns true if the strings s and t, interpreted as UTF-8 strings, are equal under Unicode case-folding.",
				Value:       equalFold,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The first string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "t",
						Description: "The second string to check.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "True if the strings are equal under Unicode case-folding.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "index",
				Description: "Returns the index of the first instance of substr in s, or -1 if substr is not present in s.",
				Value:       index,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "substr",
						Description: "The substring to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "i",
						Description: "The index of the first instance of substr in s, or -1 if substr is not present in s.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "index_any",
				Description: "Returns the index of the first instance of any of the runes in chars in s, or -1 if none of the runes in chars are present in s.",
				Value:       indexAny,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "chars",
						Description: "The characters to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "i",
						Description: "The index of the first instance of any of the runes in chars in s, or -1 if none of the runes in chars are present in s.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "last_index",
				Description: "Returns the index of the last instance of substr in s, or -1 if substr is not present in s.",
				Value:       lastIndex,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "substr",
						Description: "The substring to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "i",
						Description: "The index of the last instance of substr in s, or -1 if substr is not present in s.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "last_index_any",
				Description: "Returns the index of the last instance of any of the runes in chars in s, or -1 if none of the runes in chars are present in s.",
				Value:       lastIndexAny,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to check.",
						Type:        luadoc.String,
					},
					{
						Name:        "chars",
						Description: "The characters to check for.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "i",
						Description: "The index of the last instance of any of the runes in chars in s, or -1 if none of the runes in chars are present in s.",
						Type:        luadoc.Number,
					},
				},
			},
			{
				Name:        "duplicate",
				Description: `Returns a new string consisting of count copies of the string s. If count is zero or negative, returns the empty string.`,
				Value:       duplicate,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to repeat.",
						Type:        luadoc.String,
					},
					{
						Name:        "count",
						Description: "The number of times to repeat the string.",
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "s",
						Description: "A new string consisting of count copies of the string s.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "title",
				Description: "Returns a copy of the string s with all Unicode letters that begin words mapped to their title case.",
				Value:       title,
				Params: []*luadoc.Param{
					{
						Name:        "s",
						Description: "The string to convert.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "s",
						Description: "A copy of the string s with all Unicode letters that begin words mapped to their title case.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func replaceAll(L *lua.LState) int {
	s := L.CheckString(1)
	old := L.CheckString(2)
	neww := L.CheckString(3)

	L.Push(lua.LString(strings.ReplaceAll(s, old, neww)))
	return 1
}

func replace(L *lua.LState) int {
	s := L.CheckString(1)
	old := L.CheckString(2)
	neww := L.CheckString(3)
	n := L.CheckInt(4)

	L.Push(lua.LString(strings.Replace(s, old, neww, n)))
	return 1
}

func split(L *lua.LState) int {
	s := L.CheckString(1)
	sep := L.CheckString(2)

	table, _ := luautil.ToLValue(L, strings.Split(s, sep))
	L.Push(table)
	return 1
}

func toLower(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.ToLower(s)))
	return 1
}

func toUpper(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.ToUpper(s)))
	return 1
}

func trimSpace(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.TrimSpace(s)))
	return 1
}

func trimPrefix(L *lua.LState) int {
	s := L.CheckString(1)
	prefix := L.CheckString(2)

	L.Push(lua.LString(strings.TrimPrefix(s, prefix)))
	return 1
}

func trimSuffix(L *lua.LState) int {
	s := L.CheckString(1)
	suffix := L.CheckString(2)

	L.Push(lua.LString(strings.TrimSuffix(s, suffix)))
	return 1
}

func trim(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LString(strings.Trim(s, cutset)))
	return 1
}

func trimLeft(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LString(strings.TrimLeft(s, cutset)))
	return 1
}

func trimRight(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LString(strings.TrimRight(s, cutset)))
	return 1
}

func hasPrefix(L *lua.LState) int {
	s := L.CheckString(1)
	prefix := L.CheckString(2)

	L.Push(lua.LBool(strings.HasPrefix(s, prefix)))
	return 1
}

func hasSuffix(L *lua.LState) int {
	s := L.CheckString(1)
	suffix := L.CheckString(2)

	L.Push(lua.LBool(strings.HasSuffix(s, suffix)))
	return 1
}

func contains(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LBool(strings.Contains(s, sub)))
	return 1
}

func containsAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LBool(strings.ContainsAny(s, cutset)))
	return 1
}

func count(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LNumber(strings.Count(s, sub)))
	return 1
}

func countAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LNumber(strings.Count(s, cutset)))
	return 1
}

func equalFold(L *lua.LState) int {
	s := L.CheckString(1)
	t := L.CheckString(2)

	L.Push(lua.LBool(strings.EqualFold(s, t)))
	return 1
}

func index(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LNumber(strings.Index(s, sub)))
	return 1
}

func indexAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LNumber(strings.IndexAny(s, cutset)))
	return 1
}

func lastIndex(L *lua.LState) int {
	s := L.CheckString(1)
	sub := L.CheckString(2)

	L.Push(lua.LNumber(strings.LastIndex(s, sub)))
	return 1
}

func lastIndexAny(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)

	L.Push(lua.LNumber(strings.LastIndexAny(s, cutset)))
	return 1
}

func duplicate(L *lua.LState) int {
	s := L.CheckString(1)
	count := L.CheckInt(2)

	L.Push(lua.LString(strings.Repeat(s, count)))
	return 1
}

func title(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(strings.ToTitle(s)))
	return 1
}
