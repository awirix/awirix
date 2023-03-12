package main

import (
	"fmt"
	"os"

	"github.com/awirix/awirix/cmd"
	"github.com/awirix/awirix/config"
	"github.com/awirix/awirix/log"
	"github.com/samber/lo"
)

func handlePanic() {
	if err := recover(); err != nil {
		log.WriteErrorf(os.Stderr, "program crashed")
		fmt.Fprintf(os.Stderr, "\n\n%s", err)
		os.Exit(1)
	}
}

func init() {
	defer handlePanic()

	lo.Must0(config.Init())
	lo.Must0(log.Init())
}

func main() {
	defer handlePanic()

	cmd.Execute()
}
