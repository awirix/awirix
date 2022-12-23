package main

import (
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/cmd"
	"github.com/vivi-app/vivi/config"
	"github.com/vivi-app/vivi/log"
)

func init() {
	lo.Must0(config.Init())
	lo.Must0(log.Init())
}

func main() {
	cmd.Execute()
}
