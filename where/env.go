package where

import (
	"github.com/vivi-app/vivi/app"
	"strings"
)

// EnvConfigPath is the environment variable name for the config path
var EnvConfigPath = strings.ToUpper(app.Name) + "_CONFIG_PATH"
