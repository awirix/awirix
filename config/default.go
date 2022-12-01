package config

import (
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/key"
)

// fields is the config fields with their default values and descriptions
var fields = []*Field{
	// LOGS
	{
		key.LogsWrite,
		false,
		"Write logs to file",
	},
	{
		key.LogsLevel,
		"info",
		`Logs level.
Available options are: (from less to most verbose)
panic, fatal, error, warn, info, debug, trace`,
	},
	// END LOGS

	// PATH
	{
		key.PathAnimeDownloads,
		".",
		"Path to anime downloads",
	},
	{
		key.PathMoviesDownloads,
		".",
		"Path to movies downloads",
	},
	{
		key.PathShowsDownloads,
		".",
		"Path to shows downloads",
	},
	// END PATH

	// VIDEO
	{
		key.VideoDefaultPlayer,
		"auto",
		`Default video player.
'auto' is a special value that will try to use the most suitable player.`,
	},
	// END VIDEO

	// EXTENSIONS
	{
		key.ExtensionsSafeMode,
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
