package where

import (
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/constant"
	"os"
	"path/filepath"
)

// Config path
// Will create the directory if it doesn't exist
func Config() string {
	var path string

	if customDir, present := os.LookupEnv(EnvConfigPath); present {
		path = customDir
	} else {
		path = filepath.Join(lo.Must(os.UserConfigDir()), constant.App)
	}

	return mkdir(path)
}

// Logs path
// Will create the directory if it doesn't exist
func Logs() string {
	return mkdir(filepath.Join(Config(), "logs"))
}

// Cache path
// Will create the directory if it doesn't exist
func Cache() string {
	genericCacheDir, err := os.UserCacheDir()
	if err != nil {
		genericCacheDir = "."
	}

	cacheDir := filepath.Join(genericCacheDir, constant.PrefixCache)
	return mkdir(cacheDir)
}

// Temp path
// Will create the directory if it doesn't exist
func Temp() string {
	tempDir := filepath.Join(os.TempDir(), constant.PrefixTemp)
	return mkdir(tempDir)
}

func Downloads() string {
	path := viper.GetString(constant.DownloaderPath)
	path = os.ExpandEnv(path)
	absPath, err := filepath.Abs(path)

	if err == nil {
		path = absPath
	}

	return mkdir(path)
}
