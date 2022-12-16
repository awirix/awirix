package where

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/filesystem"
	"os"
	"path/filepath"
	"strings"
)

// mkdir creates a directory and all parent directories if they don't exist
// will return the path of the directory
func mkdir(path string) string {
	path = expand(path)
	lo.Must0(filesystem.Api().MkdirAll(path, os.ModePerm))
	return path
}

func expand(path string) string {
	if strings.Contains(path, "$") {
		envs, ok := os.LookupEnv(path)
		if !ok {
			panic(fmt.Sprintf("Error while processing path with ENV variables '%s': some variables are not found", path))
		}

		path = envs
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Sprintf("Error while getting absolute path for '%s': %s", abs, err.Error()))
	}

	return abs
}
