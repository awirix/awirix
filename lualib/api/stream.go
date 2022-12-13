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

func openDefault(L *lua.LState) int {
	url := L.CheckString(1)
	app := L.Get(2)

	var err error
	if app.Type() == lua.LTNil {
		err = open.Start(url)
	} else {
		err = open.StartWith(url, app.String())
	}

	if err != nil {
		L.RaiseError(fmt.Sprintf("error while opening url: %s", err))
	}

	return 0
}

func playVideo(L *lua.LState) int {
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
			fmt.Println("player not found:", player)
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
