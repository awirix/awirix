package config

import (
	"github.com/spf13/viper"
	"github.com/awirix/awirix/key"
	"runtime"
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
		key.PathDownloads,
		".",
		"Default downloads path",
	},
	{
		key.PathExtensions,
		"",
		"Extensions path. Leave empty for default",
	},
	{
		key.PathLogs,
		"",
		"Logs path. Leave empty for default",
	},
	// END

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

	// TUI
	{
		key.TUIClickable,
		runtime.GOOS == "android",
		"Enable support for mouse clicks in TUI",
	},
	{
		key.TUIPaddingLeft,
		0,
		"Left padding for TUI",
	},
	{
		key.TUIPaddingRight,
		0,
		"Right padding for TUI",
	},
	{
		key.TUIPaddingTop,
		0,
		"Top padding for TUI",
	},
	{
		key.TUIPaddingBottom,
		0,
		"Bottom padding for TUI",
	},
	// END TUI
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
