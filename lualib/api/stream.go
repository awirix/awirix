package api

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vivi-app/libopen/open"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
	"runtime"
)

func openDefault(L *lua.LState) int {
	url := L.CheckString(1)
	app := L.CheckAny(2)

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

	if player := viper.GetString(constant.VideoDefaultPlayer); player != "auto" {
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
		if !util.ProgramInPath(player) {
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
