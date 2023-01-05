package filepath

import (
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luadoc"
	"path/filepath"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "filepath",
		Description: "Filepath manipulation functions.",
		Funcs: []*luadoc.Func{
			{
				Name:        "join",
				Description: "Joins two or more path elements into a single path, adding a Separator if necessary. The result is Cleaned; in particular, all empty strings are ignored.",
				Value:       join,
				Params: []*luadoc.Param{
					{
						Name:        "elems",
						Description: "Path elements to join.",
						Type:        luadoc.List(luadoc.String),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The joined path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "split",
				Description: "Splits path immediately following the final Separator, separating it into a directory and file name component. If there is no Separator in path, Split returns an empty dir and file set to path. The returned values have the property that path = dir+file.",
				Value:       split,
				Params: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The path to split.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "dir",
						Description: "The directory of the path.",
						Type:        luadoc.String,
					},
					{
						Name:        "file",
						Description: "The file of the path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "clean",
				Description: "Clean returns the shortest path name equivalent to path by purely lexical processing. It applies the following rules iteratively until no further processing can be done:",
				Value:       clean,
				Params: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The path to clean.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "cleaned",
						Description: "The cleaned path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "base",
				Description: "Base returns the last element of path. Trailing path separators are removed before extracting the last element. If the path is empty, Base returns \".\". If the path consists entirely of separators, Base returns a single separator.",
				Value:       base,
				Params: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The path to get the base of.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "base",
						Description: "The base of the path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "dir",
				Description: "Dir returns all but the last element of path, typically the path's directory. After dropping the final element, the path is Cleaned and trailing slashes are removed. If the path is empty, Dir returns \".\". If the path consists entirely of separators, Dir returns a single separator. The returned path does not end in a separator unless it is the root directory.",
				Value:       dir,
				Params: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The path to get the directory of.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "dir",
						Description: "The directory of the path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "ext",
				Description: "Ext returns the file name extension used by path. The extension is the suffix beginning at the final dot in the final slash-separated element of path; it is empty if there is no dot.",
				Value:       ext,
				Params: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The path to get the extension of.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ext",
						Description: "The extension of the path.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "isAbs",
				Description: "IsAbs reports whether the path is absolute.",
				Value:       isAbs,
				Params: []*luadoc.Param{
					{
						Name:        "path",
						Description: "The path to check.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the path is absolute.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "rel",
				Description: "Rel returns a relative path that is lexically equivalent to targpath when joined to basepath with an intervening separator. That is, Join(basepath, Rel(basepath, targpath)) equals targpath. The returned path is cleaned and trailing slashes are removed. If the paths share no elements, Rel returns an error.",
				Value:       rel,
				Params: []*luadoc.Param{
					{
						Name:        "basepath",
						Description: "The base path.",
						Type:        luadoc.String,
					},
					{
						Name:        "targpath",
						Description: "The target path.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "rel",
						Description: "The relative path.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func join(L *lua.LState) int {
	L.Push(lua.LString(filepath.Join(L.CheckString(1), L.CheckString(2))))
	return 1
}

func split(L *lua.LState) int {
	dir, file := filepath.Split(L.CheckString(1))
	L.Push(lua.LString(dir))
	L.Push(lua.LString(file))
	return 1
}

func clean(L *lua.LState) int {
	L.Push(lua.LString(filepath.Clean(L.CheckString(1))))
	return 1
}

func base(L *lua.LState) int {
	L.Push(lua.LString(filepath.Base(L.CheckString(1))))
	return 1
}

func dir(L *lua.LState) int {
	L.Push(lua.LString(filepath.Dir(L.CheckString(1))))
	return 1
}

func ext(L *lua.LState) int {
	L.Push(lua.LString(filepath.Ext(L.CheckString(1))))
	return 1
}

func isAbs(L *lua.LState) int {
	L.Push(lua.LBool(filepath.IsAbs(L.CheckString(1))))
	return 1
}

func rel(L *lua.LState) int {
	rel, err := filepath.Rel(L.CheckString(1), L.CheckString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(rel))
	return 1
}
