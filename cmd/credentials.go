package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-openapi/strfmt"
	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/manifoldco/grafton/config"
	"github.com/manifoldco/grafton/generated/connector/client"
	"github.com/manifoldco/grafton/generated/connector/client/o_auth"
	"github.com/manifoldco/grafton/generated/connector/models"
)

const passwordMask = '●'
const defaultHostname = "manifold.co"
const defaultScheme = "https"

func apiURLPattern() string {
	scheme := os.Getenv("MANIFOLD_SCHEME")
	if scheme == "" {
		scheme = defaultScheme
	}
	hostname := os.Getenv("MANIFOLD_HOSTNAME")
	if hostname == "" {
		hostname = defaultHostname
	}
	return fmt.Sprintf("%s://api.%s.%s/v1", scheme, "%s", hostname)
}

var credentialFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "provider",
		Usage: "The label of the provider to rotate the credentials for",
	},
}

func init() {
	cmd := &cli.Command{
		Name:  "credentials",
		Usage: "Manage OAuth 2 credential pairs for Manifold.co",
		Subcommands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "List all existing credentials for a provider",
				Flags:  credentialFlags,
				Action: listCredentialsCmd,
			},
			{
				Name:   "rotate",
				Usage:  "Creates a new credential for a provider and sets the old one to expire in 24h",
				Flags:  credentialFlags,
				Action: createCredentialsCmd,
			},
			{
				Name:      "delete",
				ArgsUsage: "id",
				Usage:     "Delete a credential",
				Action:    deleteCredentialsCmd,
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

	provider, err := findProvider(ctx, cliCtx, client)
	if err != nil {
		return cli.NewExitError("Failed to find product: "+err.Error(), -1)
	}

	connector, err := NewConnector(token)
	if err != nil {
		return cli.NewExitError("Failed to create connector client: "+err.Error(), -1)
	}

	params := o_auth.NewPostCredentialsParamsWithContext(ctx)
	desc := "grafton rotation"

	body := &models.OAuthCredentialCreateRequest{
		Description: &desc,
		ProviderID:  &provider.ID,
	}
	params.SetBody(body)

	res, err := connector.OAuth.PostCredentials(params, nil)
	if err != nil {
		return cli.NewExitError("Failed to rotate credentials: "+err.Error(), -1)
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

	provider, err := findProvider(ctx, cliCtx, client)
	if err != nil {
		return cli.NewExitError("Failed to find provider "+err.Error(), -1)
	}

	connector, err := NewConnector(token)
	if err != nil {
		return cli.NewExitError("Failed to create connector client "+err.Error(), -1)
	}

	params := o_auth.NewGetCredentialsParamsWithContext(ctx)
	pID := provider.ID.String()
	params.SetProviderID(&pID)

	res, err := connector.OAuth.GetCredentials(params, nil)
	if err != nil {
		return cli.NewExitError("Failed to get credentials "+err.Error(), -1)
	}

	payload := res.Payload

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)

	fmt.Fprintln(w, "ID\tCreated\tExpires")

	for _, cred := range payload {
		date := time.Time(cred.ExpiresAt)
		expires := "-"

		if !date.IsZero() {
			expires = date.Format("2006-01-02 15:04:05 MST")
		}

		created := time.Time(*cred.CreatedAt).Format("2006-01-02 15:04:05 MST")
		fmt.Fprintf(w, "%s\t%s\t%s\n", cred.ID, created, expires)
	}

	w.Flush()

	return nil
}

func deleteCredentialsCmd(cliCtx *cli.Context) error {
	ctx := context.Background()

	args := cliCtx.Args()

	if args.Len() != 1 {
		cli.ShowCommandHelpAndExit(cliCtx, cliCtx.Command.Name, -1)
		return nil
	}

	id := args.First()

	_, token, err := login(ctx)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	params := o_auth.NewDeleteCredentialsIDParamsWithContext(ctx)
	params.SetID(id)

	connector, err := NewConnector(token)
	if err != nil {
		return cli.NewExitError("Failed to create connector client: "+err.Error(), -1)
	}

	_, err = connector.OAuth.DeleteCredentialsID(params, nil)
	if err != nil {
		return cli.NewExitError("Failed to delete credential "+err.Error(), -1)
	}

	fmt.Println("Credential deleted!")

	return nil
}

// NewConnector creates a new connector client with the provided 'token'
func NewConnector(token string) (*client.Connector, error) {
	u, err := url.Parse(fmt.Sprintf(apiURLPattern(), "connector"))
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

	cfgs = append(cfgs, manifold.ForURLPattern(apiURLPattern()))
	cfgs = append(cfgs, manifold.WithUserAgent("grafton-"+config.Version))

	client := manifold.New(cfgs...)

	token, err := client.Login(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	cfgs = append(cfgs, manifold.WithAPIToken(token))
	client = manifold.New(cfgs...)

	return client, token, nil
}

func findProvider(ctx context.Context, cliCtx *cli.Context, client *manifold.Client) (*manifold.Provider, error) {
	providerLabel := cliCtx.String("provider")
	if providerLabel == "" {
		return nil, errors.New("--provider flag missing")
	}

	var provider *manifold.Provider

	provList := client.Providers.List(ctx, nil)

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

	return provider, nil
}
