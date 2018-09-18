package main

import (
	"encoding/base64"
	"fmt"

	"github.com/urfave/cli"

	"github.com/manifoldco/grafton/connector"
)

func init() {
	cmd := cli.Command{
		Name:   "serve",
		Usage:  "Serves a local version of of the Connector API",
		Action: serveCmd,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "product",
				Usage:  "The label of the product being provisioned",
				EnvVar: "PRODUCT",
			},
			cli.StringFlag{
				Name:   "client-id",
				Usage:  "Client ID to use for SSO and local Connector API testing",
				EnvVar: "OAUTH2_CLIENT_ID",
			},
			cli.StringFlag{
				Name:   "client-secret",
				Usage:  "Client secret to use for SSO and local Connector API testing",
				EnvVar: "OAUTH2_CLIENT_SECRET",
			},
			cli.UintFlag{
				Name:   "connector-port",
				Usage:  "Local port for running the fake Connector API for SSO and Async testing",
				EnvVar: "CONNECTOR_PORT",
			},
		},
	}

	cmds = append(cmds, cmd)
}

func serveCmd(ctx *cli.Context) error {
	product := ctx.String("product")
	if product == "" {
		return cli.NewExitError("The 'product' flag is required and was not provided", -1)
	}
	clientID := ctx.String("client-id")
	clientSecret := ctx.String("client-secret")
	connectorPort := ctx.Uint("connector-port")
	if connectorPort == 0 {
		fmt.Println("'connector-port' was not defined, using: 3001")
		connectorPort = 3001
	}

	if clientID == "" && clientSecret == "" {
		// Attempt to infer from keyFile
		keyFile, err := getKeyFilePath()
		if err != nil {
			return cli.NewExitError("Could not determine working directory, "+
				"while attempting to infer client-id and secret: "+err.Error(), -1)
		}

		kp, err := loadKeypair(keyFile)
		if err != nil {
			return cli.NewExitError("Error reading the keypair file, "+
				"while attempting to infer client-id and secret: "+err.Error(), -1)
		}

		fmt.Println("Inferred OAuth keys from key file:")
		clientID = base64.StdEncoding.EncodeToString(kp.PublicKey)
		fmt.Println("'client-id': " + clientID)
		clientSecret = base64.StdEncoding.EncodeToString(kp.PrivateKey)
		fmt.Println("'client-secret': " + clientSecret)
	} else if clientID == "" {
		return cli.NewExitError("The 'client-id' flag is required and was not provided", -1)
	} else if clientSecret == "" {
		return cli.NewExitError("The 'client-secret' flag is required and was not provided", -1)
	}

	fakeConnector, err := connector.New(connectorPort, clientID, clientSecret, product)
	if err != nil {
		return cli.NewExitError("Error while configuring connector service: "+err.Error(), -1)
	}

	fmt.Printf("Starting Connector server on localhost:%d", connectorPort)
	return fakeConnector.StartSync()
}
