package manager

import (
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
)

func Remove(ext *extension.Extension) error {
	return filesystem.Api().RemoveAll(ext.Path())
}
