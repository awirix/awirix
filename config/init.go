package config

import (
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/where"
	"github.com/spf13/viper"
	"strings"
)

var EnvKeyReplacer = strings.NewReplacer(".", "_")

func Init() error {
	viper.SetConfigName(constant.App)
	viper.SetConfigType(constant.ConfigFormat)
	viper.SetFs(filesystem.Api())
	viper.AddConfigPath(where.Config())
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvPrefix(constant.App)
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
