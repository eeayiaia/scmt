package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func Start(defaultFn func(*cli.Context)) {
	app := cli.NewApp()
	app.Name = "scmt"
	app.Usage = "SuperK Cluster Management Toolkit"
	app.Commands = getCommands()
	app.Action = defaultFn
	app.Run(os.Args)
}
