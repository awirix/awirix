package manager

import (
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
)

func Remove(ext *extension.Extension) error {
	// remove extension from favorites on removal
	if IsFavorite(ext) {
		_ = ToggleFavorite(ext)
	}

	return filesystem.Api().RemoveAll(ext.Path())
}
