package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/asaskevich/govalidator"
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
const apiURL = "http://api.%s.manifold.co/v1"

func init() {
	cmd := cli.Command{
		Name:  "credentials",
		Usage: "Rotates OAuth credentials for Manifold.co",
		Subcommands: []cli.Command{
			{
				Name:   "rotate",
				Usage:  "Creates a new credential and set the old one to expire in 24h",
				Action: createCredentialsCmd,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "provider",
						Usage: "The label of the provider",
					},
					cli.StringFlag{
						Name:  "product",
						Usage: "The label of the product",
					},
				},
			},
		},
	}

	cmds = append(cmds, cmd)
}

func createCredentialsCmd(cliCtx *cli.Context) error {
	ctx := context.Background()

	providerLabel := cliCtx.String("provider")
	if providerLabel == "" {
		return cli.NewExitError("--provider flag missing", -1)
	}

	productLabel := cliCtx.String("product")
	if productLabel == "" {
		return cli.NewExitError("--product flag missing", -1)
	}

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
		return cli.NewExitError(err.Error(), -1)
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
		return cli.NewExitError(err.Error(), -1)
	}
	cfgs := []manifold.ConfigFunc{}

	cfgs = append(cfgs, manifold.ForURLPattern(apiURL))
	cfgs = append(cfgs, manifold.WithUserAgent("grafton-"+version))

	client := manifold.New(cfgs...)

	token, err := client.Login(ctx, email, password)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	cfgs = append(cfgs, manifold.WithAPIToken(token))
	client = manifold.New(cfgs...)

	var provider *manifold.Provider
	var product *manifold.Product

	provList := client.Providers.List(ctx)

	for provList.Next() {
		p, err := provList.Current()
		if err != nil {
			return cli.NewExitError("Fetching provider error: "+err.Error(), -1)
		}

		if p.Body.Label == providerLabel {
			provider = p
			break
		}
	}

	if provider == nil {
		return cli.NewExitError(fmt.Sprintf("Provider %q not found", providerLabel), -1)
	}

	prodList := client.Products.List(ctx, nil)

	for prodList.Next() {
		p, err := prodList.Current()
		if err != nil {
			return cli.NewExitError("Fetching product error: "+err.Error(), -1)
		}

		if p.Body.Label == productLabel {
			product = p
			break
		}
	}

	if product == nil {
		return cli.NewExitError(fmt.Sprintf("product %q not found", productLabel), -1)
	}

	connector, err := NewConnector(token)
	if err != nil {
		return cli.NewExitError("Failed to create connector client "+err.Error(), -1)
	}

	desc := "grafton rotation"

	params := o_auth.NewPostCredentialsParamsWithContext(ctx)
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
