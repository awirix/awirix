package config

import (
	"github.com/spf13/viper"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/where"
	"strings"
)

var (
	EnvKeyReplacer = strings.NewReplacer(".", "_")
	Format         = "toml"
)

func Init() error {
	viper.SetConfigName(app.Name)
	viper.SetConfigType(Format)
	viper.SetFs(filesystem.Api())
	viper.AddConfigPath(where.Config())
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvPrefix(app.Name)
	viper.SetEnvKeyReplacer(EnvKeyReplacer)

	setDefaults()

	err := viper.ReadInConfig()

	switch err.(type) {
	case viper.ConfigFileNotFoundError:
		// Use defaults then
		return nil
	default:
		return err
	}
}
