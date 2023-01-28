package where

import (
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/key"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func Home() string {
	home, _ := lo.TryOr(os.UserHomeDir, ".")
	return home
}

// Config path
// Will create the directory if it doesn't exist
func Config() string {
	var path string

	if customDir, present := os.LookupEnv(EnvConfigPath); present {
		path = customDir
	} else {
		osConfigDir, _ := lo.TryOr(os.UserConfigDir, filepath.Join(Home(), "config"))
		path = filepath.Join(osConfigDir, app.Name)
	}

	return mkdir(path)
}

// Logs path
// Will create the directory if it doesn't exist
func Logs() string {
	path := viper.GetString(key.PathLogs)
	if len(path) == 0 {
		path = filepath.Join(Cache(), "logs")
	}

	return mkdir(path)
}

// Cache path
// Will create the directory if it doesn't exist
func Cache() string {
	osCacheDir, _ := lo.TryOr(os.UserCacheDir, filepath.Join(Config(), "cache"))

	cacheDir := filepath.Join(osCacheDir, app.Name)
	return mkdir(cacheDir)
}

// Temp path
// Will create the directory if it doesn't exist
func Temp() string {
	tempDir := filepath.Join(os.TempDir(), app.Name)
	return mkdir(tempDir)
}

func Downloads() string {
	path := viper.GetString(key.PathDownloads)
	return mkdir(path)
}

func Extensions() string {
	path := viper.GetString(key.PathExtensions)

	if path == "" {
		path = filepath.Join(Config(), "extensions")
	}

	return mkdir(path)
}
