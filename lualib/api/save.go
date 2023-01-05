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
	path := L.CheckString(2)

	if path == "" {
		L.ArgError(2, "path must not be empty")
		return 0
	}

	path = filepath.Join(ext.Downloads(), path)

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
