package main

import (
	"github.com/vivi-app/vivi/cmd"
	"github.com/vivi-app/vivi/config"
	"github.com/vivi-app/vivi/log"
	"github.com/samber/lo"
)

func main() {
	// prepare config and logs
	lo.Must0(config.Init())
	lo.Must0(log.Init())

	// run the app
	cmd.Execute()
}
