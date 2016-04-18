package cli

import (
	"github.com/codegangsta/cli"
	"os"
)

func Start() {
	app := cli.NewApp()
	app.Name = "scmt"
	app.Usage = "SuperK Cluster Management Toolkit"
	app.Commands = getCommands()
	app.Run(os.Args)
}
