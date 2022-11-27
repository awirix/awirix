package config

import (
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/constant"
)

// fields is the config fields with their default values and descriptions
var fields = []*Field{
	// LOGS
	{
		constant.LogsWrite,
		false,
		"Write logs to file",
	},
	{
		constant.LogsLevel,
		"info",
		`Logs level.
Available options are: (from less to most verbose)
panic, fatal, error, warn, info, debug, trace`,
	},
	// END LOGS

	// PATH
	{
		constant.PathAnimeDownloads,
		".",
		"Path to anime downloads",
	},
	{
		constant.PathMoviesDownloads,
		".",
		"Path to movies downloads",
	},
	{
		constant.PathShowsDownloads,
		".",
		"Path to shows downloads",
	},
	// END PATH

	// VIDEO
	{
		constant.VideoDefaultPlayer,
		"auto",
		`Default video player.
'auto' is a special value that will try to use the most suitable player.`,
	},
	// END VIDEO

	// EXTENSIONS
	{
		constant.ExtensionsSafeMode,
		true,
		`Enable safe mode for extensions.
If enabled, system commands will be disabled
and the extension will be unable to access the filesystem.`,
	},
	// EXTENSIONS
}

func setDefaults() {
	Default = make(map[string]*Field, len(fields))
	for _, f := range fields {
		Default[f.Key] = f
		viper.SetDefault(f.Key, f.DefaultValue)
		viper.MustBindEnv(f.Key)
	}
}

var Default map[string]*Field
