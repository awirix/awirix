package manager

import (
	"fmt"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/where"
	"github.com/samber/mo"
	"path/filepath"
)

var installed = mo.None[[]*extension.Extension]()

func ResetInstalledCache() {
	installed = mo.None[[]*extension.Extension]()
}

func Installed() ([]*extension.Extension, error) {
	if exts, ok := installed.Get(); ok {
		return exts, nil
	}

	path := where.Extensions()

	dir, err := filesystem.Api().ReadDir(path)
	if err != nil {
		return nil, err
	}

	extensions := make([]*extension.Extension, 0)

	for _, ext := range dir {
		if !ext.IsDir() {
			continue
		}

		path := filepath.Join(path, ext.Name())
		ext, err := extension.New(path)
		if err != nil {
			log.Errorf("failed to load extension at '%s': %s", path, err)
			continue
		}

		extensions = append(extensions, ext)
	}

	installed = mo.Some(extensions)
	return extensions, nil
}

func GetExtensionByID(id string) (*extension.Extension, error) {
	extensions, err := Installed()
	if err != nil {
		return nil, err
	}

	for _, ext := range extensions {
		if ext.Passport().ID == id {
			return ext, nil
		}
	}

	return nil, fmt.Errorf("extension not found: %s", id)
}
