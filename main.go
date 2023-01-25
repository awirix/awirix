package main

import (
	"github.com/awirix/awirix/cmd"
	"github.com/awirix/awirix/config"
	"github.com/awirix/awirix/log"
	"github.com/samber/lo"
)

func init() {
	lo.Must0(config.Init())
	lo.Must0(log.Init())
}

func main() {
	cmd.Execute()
}
