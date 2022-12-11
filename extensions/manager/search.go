package manager

import (
	"fmt"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/where"
	"path/filepath"
)

func InstalledExtensions() ([]*extension.Extension, error) {
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

			ext, err := extension.New(filepath.Join(path, d.Name()))
			if err != nil {
				return nil, err
			}

			extensions = append(extensions, ext)
		}
	}

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
