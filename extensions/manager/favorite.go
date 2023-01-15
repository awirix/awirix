package manager

import (
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/where"
	"github.com/metafates/gache"
	"path/filepath"
)

var favoriteExtensions = gache.New[map[string]*extension.Extension](&gache.Options{
	Path:       filepath.Join(where.Cache(), "favorite_extensions"),
	FileSystem: &filesystem.GacheFs{},
})

func ToggleFavorite(ext *extension.Extension) error {
	cached, expired, err := favoriteExtensions.Get()
	if err != nil {
		return err
	}

	if expired || cached == nil {
		cached = make(map[string]*extension.Extension)
	}

	id := ext.Passport().ID
	if _, ok := cached[id]; ok {
		delete(cached, id)
	} else {
		cached[id] = ext
	}

	return favoriteExtensions.Set(cached)
}

func IsFavorite(ext *extension.Extension) bool {
	cached, expired, err := favoriteExtensions.Get()
	if err != nil || expired || cached == nil {
		return false
	}

	_, ok := cached[ext.Passport().ID]
	return ok
}
