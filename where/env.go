package where

import (
	"github.com/vivi-app/vivi/constant"
	"strings"
)

// EnvConfigPath is the environment variable name for the config path
var EnvConfigPath = strings.ToUpper(constant.App) + "_CONFIG_PATH"
