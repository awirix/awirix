package where

import (
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/key"
	"github.com/vivi-app/vivi/prefix"
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
		path = filepath.Join(lo.Must(os.UserConfigDir()), app.Name)
	}

	return mkdir(path)
}

// Logs path
// Will create the directory if it doesn't exist
func Logs() string {
	path := viper.GetString(key.PathLogs)
	if len(path) == 0 {
		path = filepath.Join(Config(), "logs")
	}

	return mkdir(path)
}

// Cache path
// Will create the directory if it doesn't exist
func Cache() string {
	osCacheDir, err := os.UserCacheDir()
	if err != nil {
		osCacheDir = "."
	}

	cacheDir := filepath.Join(osCacheDir, prefix.Cache)
	return mkdir(cacheDir)
}

// Temp path
// Will create the directory if it doesn't exist
func Temp() string {
	tempDir := filepath.Join(os.TempDir(), prefix.Temp)
	return mkdir(tempDir)
}

func Downloads() string {
	path := viper.GetString(key.PathDownloads)
	return mkdir(path)
}

func Extensions() string {
	path := viper.GetString(key.PathExtensions)
	if len(path) == 0 {
		path = filepath.Join(Config(), "extensions")
	}

	return mkdir(path)
}
