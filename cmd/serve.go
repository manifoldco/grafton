package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"

	"github.com/urfave/cli"

	"github.com/manifoldco/grafton/connector"
	"github.com/manifoldco/grafton/marketplace"
)

var pathRegex = regexp.MustCompile(`^(?:.*\/)?v1\/?$`)

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
			cli.StringFlag{
				Name:   "provider-api",
				Usage:  "URL for the provider API, for the Marketplace to connect with",
				EnvVar: "PROVIDER_API",
			},
			cli.UintFlag{
				Name:   "connector-port",
				Usage:  "Local port for running the fake Connector API for SSO and Async testing",
				EnvVar: "CONNECTOR_PORT",
			},
			cli.UintFlag{
				Name:   "marketplace-port",
				Usage:  "Local port for running the fake Marketplace Web Server for SSO and Async testing",
				EnvVar: "MARKETPLACE_PORT",
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
	providerAPI := ctx.String("provider-api")
	connectorPort := ctx.Uint("connector-port")
	marketplacePort := ctx.Uint("marketplace-port")
	if providerAPI == "" {
		fmt.Println("'provider-api' was not defined, using: http://localhost:3000")
		providerAPI = "http://localhost:3000/v1/"
	}
	if connectorPort == 0 {
		fmt.Println("'connector-port' was not defined, using: 3001")
		connectorPort = 3001
	}
	if marketplacePort == 0 {
		fmt.Println("'marketplace-port' was not defined, using: 3002")
		marketplacePort = 3002
	}

	pAPI, err := url.Parse(providerAPI)
	if err != nil {
		return cli.NewExitError("Failed to parse provider API URL '"+providerAPI+
			"' - "+err.Error(), -1)
	}
	if !pathRegex.Match([]byte(pAPI.Path)) {
		path := pAPI.Path
		if len(path) > 0 && pAPI.Path[len(path)-1] == '/' {
			pAPI.Path += "v1/"
		} else {
			pAPI.Path += "/v1/"
		}
		fmt.Printf("'provider-api' was missing the trailing '/v1/' specifier in the path, using: %s\n", pAPI)
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

	k, err := getKeypair()
	if err != nil {
		return err
	}
	lkp, err := k.liveKeypair()
	if err != nil {
		return cli.NewExitError("Could not create request signing keypair: "+err.Error(), -1)
	}

	fakeConnector, err := connector.New(connectorPort, clientID, clientSecret, product)
	if err != nil {
		return cli.NewExitError("Error while configuring connector service: "+err.Error(), -1)
	}
	fakeMarketplace := marketplace.New(fakeConnector, marketplacePort, pAPI, lkp)

	fmt.Printf("Starting Connector server on http://localhost:%d\n", connectorPort)
	fakeConnector.Start()
	fmt.Printf("Starting Marketplace server on http://localhost:%d\n", marketplacePort)
	return fakeMarketplace.StartSync()
}
