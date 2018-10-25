package main

import (
	"fmt"

	"github.com/manifoldco/grafton/config"
	"github.com/urfave/cli"
)

func init() {
	versionCmd := cli.Command{
		Name:   "version",
		Usage:  "Display version of utility",
		Action: versionLookup,
	}

	cmds = append(cmds, versionCmd)
}

func versionLookup(ctx *cli.Context) error {
	fmt.Printf("%s\n", config.Version)

	return nil
}
