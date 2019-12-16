package main

import (
	"os"

	"github.com/manifoldco/grafton/config"
	"github.com/urfave/cli/v2"
)

var cmds []cli.Command

func main() {
	cli.VersionPrinter = func(ctx *cli.Context) {
		versionLookup(ctx)
	}

	app := cli.NewApp()
	app.Name = "grafton"
	app.HelpName = "grafton"
	app.Usage = "Tool for testing integrations with Manifold"
	app.Version = config.Version
	app.Commands = cmds

	app.Run(os.Args)
}
