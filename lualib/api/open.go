package api

import (
	"fmt"
	"github.com/vivi-app/libopen/open"
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/where"
)

func openDefault(L *lua.LState) int {
	target := L.CheckString(1)
	app := L.OptString(2, "")

	var err error
	if app == "" {
		err = open.Start(target)
	} else {
		err = open.StartWith(target, app)
	}

	if err != nil {
		L.RaiseError(fmt.Sprintf("error while opening %s: %s", target, err))
	}

	return 0
}

func openData(L *lua.LState) int {
	ext := L.CheckString(1)
	data := L.CheckString(2)
	app := L.OptString(3, "")

	file, err := filesystem.Api().TempFile(where.Temp(), "*."+ext)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	_, err = file.WriteString(data)

	path := file.Name()

	if app == "" {
		err = open.Start(path)
	} else {
		err = open.StartWith(path, app)
	}

	return 0
}
