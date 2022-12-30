package api

import (
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/extensions"
	"github.com/vivi-app/vivi/filesystem"
	"path/filepath"
)

func save(L *lua.LState) int {
	ext := L.Context().Value("extension").(extensions.ExtensionContainer)
	data := L.CheckString(1)
	segments := L.CheckTable(2)

	path := ext.Downloads()

	if segments.Len() == 0 {
		L.ArgError(2, "table must not be empty")
		return 0
	}

	segments.ForEach(func(key lua.LValue, value lua.LValue) {
		path = filepath.Join(path, value.String())
	})

	err := filesystem.Api().MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		L.RaiseError("failed to create downloads directory: %s", err)
	}

	err = filesystem.Api().WriteFile(path, []byte(data), 0644)
	if err != nil {
		L.RaiseError("failed to save file: %s", err)
	}

	return 0
}
