package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/asaskevich/govalidator"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-openapi/strfmt"
	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/manifoldco/grafton/generated/connector/client"
	"github.com/manifoldco/grafton/generated/connector/client/o_auth"
	"github.com/manifoldco/grafton/generated/connector/models"
)

const passwordMask = '●'
const apiURL = "http://api.%s.arigato.tools/v1"

var credentialFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "provider",
		Usage: "The label of the provider",
	},
	cli.StringFlag{
		Name:  "product",
		Usage: "The label of the product",
	},
}

func init() {
	cmd := cli.Command{
		Name:  "credentials",
		Usage: "Manage OAuth 2 credential pairs for Manifold.co",
		Subcommands: []cli.Command{
			{
				Name:   "list",
				Usage:  "List all existing credentials for a product",
				Flags:  credentialFlags,
				Action: listCredentialsCmd,
			},
			{
				Name:   "rotate",
				Usage:  "Creates a new credential and sets the old one to expire in 24h",
				Flags:  credentialFlags,
				Action: createCredentialsCmd,
			},
		},
	}

	cmds = append(cmds, cmd)
}

func createCredentialsCmd(cliCtx *cli.Context) error {
	ctx := context.Background()

	client, token, err := login(ctx)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	product, err := findProduct(ctx, cliCtx, client)
	if err != nil {
		return cli.NewExitError("Failed to find product "+err.Error(), -1)
	}

	connector, err := NewConnector(token)
	if err != nil {
		return cli.NewExitError("Failed to create connector client "+err.Error(), -1)
	}

	params := o_auth.NewPostCredentialsParamsWithContext(ctx)
	desc := "grafton rotation"

	body := &models.OAuthCredentialCreateRequest{
		Description: &desc,
		ProductID:   product.ID,
	}
	params.SetBody(body)

	res, err := connector.OAuth.PostCredentials(params, nil)
	if err != nil {
		return cli.NewExitError("Failed to rotate credentials "+err.Error(), -1)
	}

	payload := res.Payload

	fmt.Println("Your old credentials will expire automatically in 24 hours")
	fmt.Println("Make sure to copy your new credentials now. You won’t be able to see this information again!")
	fmt.Printf("Client ID: %s\n", payload.ID)
	fmt.Printf("Secret: %s\n", payload.Secret)

	return nil
}

func listCredentialsCmd(cliCtx *cli.Context) error {
	ctx := context.Background()

	client, token, err := login(ctx)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	product, err := findProduct(ctx, cliCtx, client)
	if err != nil {
		return cli.NewExitError("Failed to find product "+err.Error(), -1)
	}

	connector, err := NewConnector(token)
	if err != nil {
		return cli.NewExitError("Failed to create connector client "+err.Error(), -1)
	}

	params := o_auth.NewGetCredentialsParamsWithContext(ctx)
	params.SetProductID(product.ID.String())

	res, err := connector.OAuth.GetCredentials(params, nil)
	if err != nil {
		return cli.NewExitError("Failed to get credentials "+err.Error(), -1)
	}

	payload := res.Payload

	spew.Dump(payload)

	return nil
}

func NewConnector(token string) (*client.Connector, error) {
	u, err := url.Parse(fmt.Sprintf(apiURL, "connector"))
	if err != nil {
		return nil, err
	}

	c := client.DefaultTransportConfig()
	c.WithHost(u.Host)
	c.WithBasePath(u.Path)
	c.WithSchemes([]string{u.Scheme})

	transport := httptransport.New(c.Host, c.BasePath, c.Schemes)

	transport.DefaultAuthentication = httptransport.BearerToken(token)

	return client.New(transport, strfmt.Default), nil
}

func login(ctx context.Context) (*manifold.Client, string, error) {
	fmt.Println("Please use your Manifold account to login.")
	fmt.Println("If you don't have an account yet, reach out to support@manifold.co.")

	p := promptui.Prompt{
		Label: "Email",
		Validate: func(input string) error {
			valid := govalidator.IsEmail(input)
			if valid {
				return nil
			}

			return errors.New("Please enter a valid email address")
		},
	}

	email, err := p.Run()
	if err != nil {
		return nil, "", err
	}

	p = promptui.Prompt{
		Label: "Password",
		Mask:  passwordMask,
		Validate: func(input string) error {
			if len(input) < 8 {
				return errors.New("Passwords must be greater than 8 characters")
			}

			return nil
		},
	}

	password, err := p.Run()
	if err != nil {
		return nil, "", err
	}
	cfgs := []manifold.ConfigFunc{}

	cfgs = append(cfgs, manifold.ForURLPattern(apiURL))
	cfgs = append(cfgs, manifold.WithUserAgent("grafton-"+version))

	client := manifold.New(cfgs...)

	token, err := client.Login(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	cfgs = append(cfgs, manifold.WithAPIToken(token))
	client = manifold.New(cfgs...)

	return client, token, nil
}

func findProduct(ctx context.Context, cliCtx *cli.Context, client *manifold.Client) (*manifold.Product, error) {
	providerLabel := cliCtx.String("provider")
	if providerLabel == "" {
		return nil, errors.New("--provider flag missing")
	}

	productLabel := cliCtx.String("product")
	if productLabel == "" {
		return nil, errors.New("--product flag missing")
	}

	var provider *manifold.Provider
	var product *manifold.Product

	provList := client.Providers.List(ctx)

	defer provList.Close()

	for provList.Next() {
		p, err := provList.Current()
		if err != nil {
			return nil, err
		}

		if p.Body.Label == providerLabel {
			provider = p
			break
		}
	}

	if provider == nil {
		return nil, fmt.Errorf("Provider %q not found", providerLabel)
	}

	opts := manifold.ProductsListOpts{ProviderID: &provider.ID}

	prodList := client.Products.List(ctx, &opts)
	defer prodList.Close()

	for prodList.Next() {
		p, err := prodList.Current()
		if err != nil {
			return nil, err
		}

		if p.Body.Label == productLabel {
			product = p
			break
		}
	}

	if product == nil {
		return nil, fmt.Errorf("Provider %q not found", providerLabel)
	}

	return product, nil
}
