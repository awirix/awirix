package manager

import (
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
)

func UninstallExtension(ext *extension.Extension) error {
	return filesystem.Api().RemoveAll(ext.Path())
}
