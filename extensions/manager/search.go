package manager

import (
	"fmt"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/option"
	"github.com/awirix/awirix/where"
	"path/filepath"
)

var installed = option.None[[]*extension.Extension]()

func ResetInstalledCache() {
	installed = option.None[[]*extension.Extension]()
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

	for _, owner := range dir {
		if !owner.IsDir() {
			continue
		}

		path := filepath.Join(path, owner.Name())
		dir, err := filesystem.Api().ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, d := range dir {
			if !d.IsDir() {
				continue
			}

			extensionPath := filepath.Join(path, d.Name())
			ext, err := extension.New(extensionPath)
			if err != nil {
				log.Errorf("failed to load extension at '%s': %s", extensionPath, err)
				continue
			}

			extensions = append(extensions, ext)
		}
	}

	installed = option.Some(extensions)
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
