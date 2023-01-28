package where

import (
	"fmt"
	"github.com/awirix/awirix/filesystem"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
)

// mkdir creates a directory and all parent directories if they don't exist
// will return the path of the directory
func mkdir(path string) string {
	path = ExpandPath(path)
	lo.Must0(filesystem.Api().MkdirAll(path, os.ModePerm))
	return path
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Errorf("cannot expand path %s: %w", path, err))
		}

		path = home + path[1:]
	}

	path = os.ExpandEnv(path)
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Sprintf("error while getting absolute path for '%s': %s", abs, err.Error()))
	}

	return abs
}
