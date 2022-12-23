package manager

import (
	"fmt"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/log"
	"github.com/vivi-app/vivi/option"
	"github.com/vivi-app/vivi/where"
	"path/filepath"
)

var installedExtensions = option.None[[]*extension.Extension]()

func InstalledExtensions() ([]*extension.Extension, error) {
	if exts, ok := installedExtensions.Get(); ok {
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

	installedExtensions = option.Some(extensions)
	return extensions, nil
}

func GetExtensionByID(id string) (*extension.Extension, error) {
	extensions, err := InstalledExtensions()
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
