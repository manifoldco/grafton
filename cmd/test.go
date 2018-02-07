package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	nurl "net/url"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/acceptance"

	"github.com/manifoldco/grafton/generated/connector/models"
)

var (
	bold  = promptui.Styler(promptui.FGBold)
	faint = promptui.Styler(promptui.FGFaint)
)

func init() {
	cmd := cli.Command{
		Name:      "test",
		Usage:     "Tests the API endpoints required to integrate with Manifold",
		ArgsUsage: "[url]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "product",
				Usage:  "The label of the product being provisioned",
				EnvVar: "PRODUCT",
			},
			cli.StringFlag{
				Name:   "plan",
				Usage:  "The label of the plan for the provisioning resource",
				EnvVar: "PLAN",
			},
			cli.StringFlag{
				Name:   "plan-features",
				Usage:  "A JSON object describing the selected features for the provisioning resource",
				EnvVar: "PLAN_FEATURES",
			},
			cli.StringFlag{
				Name:   "new-plan",
				Usage:  "The plan to resize the instance to from the original plan",
				EnvVar: "NEW_PLAN",
			},
			cli.StringFlag{
				Name:   "new-plan-features",
				Usage:  "A JSON object describing the selected features for the resizing from the original plan",
				EnvVar: "NEW_PLAN_FEATURES",
			},
			cli.StringFlag{
				Name:   "region",
				Usage:  "The label of the region which the resource will be provision in",
				EnvVar: "REGION",
			},
			cli.StringSliceFlag{
				Name:   "exclude",
				Usage:  "Exclude running these feature tests (and those that depend on it)",
				EnvVar: "EXCLUDE",
			},
			cli.BoolFlag{
				Name:   "no-error-cases",
				Usage:  "Skip running the error case tests",
				EnvVar: "NO_ERROR_CASES",
			},
			cli.StringFlag{
				Name:   "log",
				Usage:  "Informational logging level during tests. One of (off, info, verbose)",
				EnvVar: "LOG",
				Value:  "off",
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
			cli.StringFlag{
				Name:   "callback-timeout",
				Usage:  "duration to wait (max. 24hours) for a callback (default: 5m)",
				EnvVar: "CALLBACK_TIMEOUT",
			},
			cli.StringFlag{
				Name:  "resource-measures",
				Usage: "Optional measures map to be returned by resource measures polling",
				Value: `{"feature-a": 0, "feature-b": 1000}`,
			},
		},
		Action: testCmd,
	}

	cmds = append(cmds, cmd)
}

func testCmd(ctx *cli.Context) error {
	args := ctx.Args()

	url := "http://localhost:3000"
	plan := ctx.String("plan")
	sPlanFeatures := ctx.String("plan-features")
	newPlan := ctx.String("new-plan")
	sNewPlanFeatures := ctx.String("new-plan-features")
	product := ctx.String("product")
	region := ctx.String("region")
	excludeFeatures := ctx.StringSlice("exclude")

	clientID := ctx.String("client-id")
	clientSecret := ctx.String("client-secret")
	connectorPort := ctx.Uint("connector-port")
	callbackTimeout := ctx.String("callback-timeout")

	resourceMeasures := ctx.String("resource-measures")

	var logLevel acceptance.LogLevel
	rawLevel := ctx.String("log")
	switch acceptance.LogLevel(rawLevel) {
	case acceptance.LogOff:
		logLevel = acceptance.LogOff
	case acceptance.LogInfo:
		logLevel = acceptance.LogInfo
	case acceptance.LogVerbose:
		logLevel = acceptance.LogVerbose
		// we need to set it so the openapi runtime gets triggered to run in
		// verbose mode and actually prints http request and response data.
		os.Setenv("DEBUG", "true")
	default:
		return cli.NewExitError("invalid log value "+rawLevel, -1)
	}

	planFeatures := models.FeatureMap{}
	if sPlanFeatures != "" {
		err := json.Unmarshal([]byte(sPlanFeatures), planFeatures)
		if err != nil {
			return cli.NewExitError("The supplied plan-features does not appear to be valid JSON: "+err.Error(), -1)
		}
	}

	newPlanFeatures := models.FeatureMap{}
	if sNewPlanFeatures != "" {
		err := json.Unmarshal([]byte(sNewPlanFeatures), newPlanFeatures)
		if err != nil {
			return cli.NewExitError("The supplied new-plan-features does not appear to be valid JSON: "+err.Error(), -1)
		}
	}

	if len(args) > 0 {
		url = args[0]
	}

	purl, err := nurl.Parse(url)
	if err != nil {
		return cli.NewExitError("unable to parse url: "+url, -1)
	}

	// Always append the '/v1' to the path
	if !strings.HasSuffix(purl.Path, "/v1") {
		purl.Path = path.Join(purl.Path, "/v1")
	}

	k, err := getKeypair()
	if err != nil {
		return err
	}

	lkp, err := k.liveKeypair()
	if err != nil {
		return cli.NewExitError("Could not create request signing keypair: "+err.Error(), -1)
	}

	connectorURL := deriveConnectorURL(connectorPort)
	api := grafton.New(purl, connectorURL, lkp, nil)

	fkp, err := emptyKeypair()
	if err != nil {
		return cli.NewExitError("Could not create request empty signing keypair: "+err.Error(), -1)
	}
	unauthorizedAPI := grafton.New(purl, connectorURL, fkp, nil)
	c := context.Background()

	willChangePlan := false
	if newPlan != "" {
		willChangePlan = true
	}

	acceptance.SetLogLevel(logLevel)

	acceptance.Infoln(bold("Configration"))
	buf := bytes.NewBufferString("")
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\tURL:\t%s\n", faint(url))
	fmt.Fprintf(w, "\tProduct:\t%s\n", faint(product))
	fmt.Fprintf(w, "\tPlan:\t%s\n", faint(plan))
	fmt.Fprintf(w, "\tRegion:\t%s\n", faint(region))
	fmt.Fprintf(w, "\tResizing?\t%s\n", faint(yn(willChangePlan)))

	if willChangePlan {
		fmt.Fprintf(w, "\tNew Plan:\t%s\n", faint(newPlan))
	}

	if len(excludeFeatures) > 0 {
		fmt.Fprintf(w, "\tExcluded Features:\t%s\n", faint(strings.Join(excludeFeatures, " ")))
	}

	fmt.Fprintf(w, "\tClient ID:\t%s\n", faint(clientID))
	fmt.Fprintf(w, "\tClient Secret:\t%s\n", faint(clientSecret))
	fmt.Fprintf(w, "\tConnector Port:\t%s\n", faint(fmt.Sprintf("%d", connectorPort)))

	if errs := acceptance.Validate(c, ctx, excludeFeatures); len(errs) != 0 {
		// format errors into a single string
		errString := []string{}
		for _, err := range errs {
			errString = append(errString, err.Error())
		}
		return cli.NewExitError(strings.Join(errString, "\n"), -1)
	}

	w.Flush()

	acceptance.Infoln(buf.String())

	cfg := acceptance.Configration{
		API:              api,
		UnauthorizedAPI:  unauthorizedAPI,
		Product:          product,
		Region:           region,
		Plan:             plan,
		PlanFeatures:     planFeatures,
		NewPlan:          newPlan,
		NewPlanFeatures:  newPlanFeatures,
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		Port:             connectorPort,
		CallbackTimeout:  callbackTimeout,
		ResourceMeasures: resourceMeasures,
	}

	if err := acceptance.Configure(cfg); err != nil {
		return cli.NewExitError("Error: "+err.Error(), -1)
	}
	if ok := acceptance.Run(c, !ctx.Bool("no-error-cases"), excludeFeatures); !ok {
		os.Exit(1)
	}

	return nil
}

func deriveConnectorURL(port uint) *nurl.URL {
	return &nurl.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d", port),
		Path:   "/v1",
	}
}

func getKeypair() (*keypair, error) {
	keyFile, err := getKeyFilePath()
	if err != nil {
		return nil, cli.NewExitError("Could not determine working directory: "+err.Error(), -1)
	}

	if _, err = os.Stat(keyFile); os.IsNotExist(err) {
		return nil, cli.NewExitError(
			"Master key file does not exist; generate one using 'grafton generate'", -1)
	}

	k, err := loadKeypair(keyFile)
	if err != nil {
		return nil, cli.NewExitError("Could not load master key file: "+err.Error(), -1)
	}

	return k, err
}

func yn(v bool) string {
	if v {
		return "yes"
	}

	return "no"
}
