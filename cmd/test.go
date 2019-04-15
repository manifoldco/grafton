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

	"github.com/urfave/cli"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/promptui"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/acceptance"
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
			cli.StringFlag{
				Name:   "import-code",
				Usage:  "The import code to import an existing resource for that resource",
				EnvVar: "IMPORT_CODE",
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
				Usage: "Optional measures map to be returned by resource measures",
				Value: `{"feature-a": 0, "feature-b": 1000}`,
			},
			cli.StringFlag{
				Name:  "credential",
				Usage: "Describes the credential type that is supported by this product. One of (single, multiple)",
				Value: "multiple",
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

	credential := ctx.String("credential")

	var logLevel acceptance.LogLevel
	rawLevel := ctx.String("log")
	switch acceptance.LogLevel(rawLevel) {
	case acceptance.LogOff:
		logLevel = acceptance.LogOff
	case acceptance.LogInfo:
		logLevel = acceptance.LogInfo
	case acceptance.LogVerbose:
		logLevel = acceptance.LogVerbose
	default:
		return cli.NewExitError("invalid log value "+rawLevel, -1)
	}

	planFeatures := manifold.FeatureMap{}
	if sPlanFeatures != "" {
		err := json.Unmarshal([]byte(sPlanFeatures), &planFeatures)
		if err != nil {
			return cli.NewExitError("The supplied plan-features does not appear to be valid JSON: "+err.Error(), -1)
		}
	}

	newPlanFeatures := manifold.FeatureMap{}
	if sNewPlanFeatures != "" {
		err := json.Unmarshal([]byte(sNewPlanFeatures), &newPlanFeatures)
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

	opt := grafton.ClientOptions{
		URL:          purl,
		ConnectorURL: deriveConnectorURL(connectorPort),
		Signer:       lkp,
		Debug:        logLevel == acceptance.LogVerbose,
	}

	api := grafton.NewClient(opt)

	fkp, err := emptyKeypair()
	if err != nil {
		return cli.NewExitError("Could not create request empty signing keypair: "+err.Error(), -1)
	}

	opt.Signer = fkp

	unauthorizedAPI := grafton.NewClient(opt)
	c := context.Background()

	willChangePlan := false
	if newPlan != "" {
		willChangePlan = true
	}

	acceptance.SetLogLevel(logLevel)

	acceptance.Infoln(bold("Configuration"))
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

	if !contains(excludeFeatures, "resource-measures") {
		fmt.Fprintf(w, "\tResource Measures:\t%s\n", faint(resourceMeasures))
	}

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

	cfg := acceptance.Configuration{
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
		Credential:       credential,
	}

	if err := acceptance.Configure(cfg); err != nil {
		return cli.NewExitError("Error: "+err.Error(), -1)
	}

	failed := acceptance.Run(c, !ctx.Bool("no-error-cases"), excludeFeatures)
	if failed {
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

func yn(v bool) string {
	if v {
		return "yes"
	}

	return "no"
}

func contains(list []string, s string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}
