package api

import (
	"fmt"
	"github.com/awirix/awirix/executil"
	"github.com/awirix/awirix/key"
	"github.com/awirix/libopen/open"
	lua "github.com/awirix/lua"
	"github.com/spf13/viper"
)

// videoPlayers the default video players that support streaming via url
var videoPlayers = [...]string{
	"iina", // macOS only
	"mpv",
	"vlc",
}

func watch(L *lua.LState) int {
	url := L.CheckString(1)

	if player := viper.GetString(key.VideoDefaultPlayer); player != "auto" {
		err := open.StartWith(url, player)

		if err != nil {
			L.RaiseError(fmt.Sprintf("error while opening video: %s", err))
		}

		return 0
	}

	for _, player := range videoPlayers {
		if !executil.ProgramInPath(player) {
			continue
		}

		err := open.StartWith(url, player)
		if err == nil {
			return 0
		} else {
			L.RaiseError(fmt.Sprintf("error while opening video: %s", err))
		}
	}

	err := open.Start(url)
	if err != nil {
		L.RaiseError(fmt.Sprintf("error while opening video: %s", err))
	}

	return 0
}
