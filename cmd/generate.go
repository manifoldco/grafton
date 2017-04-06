package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func init() {
	cmd := cli.Command{
		Name:   "generate",
		Usage:  "Generates public and private signing keys for testing Manifold API integrations locally",
		Action: generateCmd,
	}

	cmds = append(cmds, cmd)
}

func generateCmd(ctx *cli.Context) error {
	keyFile, err := getKeyFilePath()
	if err != nil {
		return cli.NewExitError("Could not determine working directory: "+err.Error(), -1)
	}

	fmt.Println("Generating Master Keypair")
	k, err := newKeypair()
	if err != nil {
		return cli.NewExitError("Could not generate keypair: "+err.Error(), -1)
	}

	fmt.Printf("Writing master keypair to file: %s\n", keyFile)
	err = k.save(keyFile)
	if err != nil {
		return cli.NewExitError("Could not write to file: "+err.Error(), -1)
	}

	fmt.Println("Success.")
	return nil
}
