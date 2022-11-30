package manager

import (
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/filesystem"
)

func UninstallExtension(ext *extension.Extension) error {
	return filesystem.Api().RemoveAll(ext.Path())
}
