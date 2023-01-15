package config

import (
	"fmt"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/key"
	"github.com/spf13/viper"
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
	{
		key.ExtensionsNewInitGitRepo,
		true,
		"Initialize a git repository when creating a new extension",
	},
	{
		key.ExtensionsNewAddLibraryDoc,
		true,
		fmt.Sprintf(`Add library documentation "%s.lua" when creating a new extension`, app.Name),
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
		1,
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
	{
		key.TUIShowExtensionAuthor,
		false,
		"Show extension's author in TUI",
	},
	{
		key.TUIShowDescription,
		true,
		"Show item's description in TUI",
	},
	{
		key.TUIPromptSymbol,
		"> ",
		"Prompt symbol to use in text input",
	},
	{
		key.TUIShowAppVersion,
		true,
		"Show app version on the top left corner of the TUI",
	},
	// END TUI

	// ICON
	{
		key.IconShowExtensionIcon,
		true,
		"Show extension's icon",
	},
	{
		key.IconShowExtensionFlag,
		false,
		"Show extension's flag based on its language",
	},
	// END ICON
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
