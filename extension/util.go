package extension

import (
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/where"
	"path/filepath"
)

func ListInstalled() []*Extension {
	extensions := make([]*Extension, 0)
	installed, err := filesystem.Api().ReadDir(where.Extensions())
	if err != nil {
		return extensions
	}

	for _, file := range installed {
		extension, err := NewFromPath(filepath.Join(where.Extensions(), file.Name()))
		if err != nil {
			continue
		}

		extensions = append(extensions, extension)
	}

	return extensions
}
