package api

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vivi-app/libopen/open"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/executil"
	"github.com/vivi-app/vivi/key"
	"runtime"
)

func watch(L *lua.LState) int {
	url := L.CheckString(1)

	if player := viper.GetString(key.VideoDefaultPlayer); player != "auto" {
		err := open.StartWith(url, player)

		if err != nil {
			L.RaiseError(fmt.Sprintf("error while opening video: %s", err))
		}

		return 0
	}

	var players = make([]string, 0)

	if runtime.GOOS == "darwin" {
		players = append(players, "iina")
	}

	players = append(players, "mpv", "vlc")

	for _, player := range players {
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
