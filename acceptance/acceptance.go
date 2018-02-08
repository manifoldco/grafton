package acceptance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/onsi/gomega"
	"github.com/urfave/cli"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/connector"

	"github.com/manifoldco/grafton/generated/connector/models"
)

var api *grafton.Client
var uapi *grafton.Client

var product string
var plan string
var planFeatures models.FeatureMap
var region string
var newPlan string
var newPlanFeatures models.FeatureMap
var resourceMeasures map[string]int64

var clientID string
var clientSecret string
var connectorPort uint
var fakeConnector *connector.FakeConnector
var cbTimeout = time.Minute * 5
var maxTimeout = 24 * time.Hour

type visitorFunc func(context.Context, *FeatureImpl) bool

type Configuration struct {
	API              *grafton.Client
	UnauthorizedAPI  *grafton.Client
	Product          string
	Region           string
	Plan             string
	PlanFeatures     models.FeatureMap
	NewPlan          string
	NewPlanFeatures  models.FeatureMap
	ClientID         string
	ClientSecret     string
	Port             uint
	CallbackTimeout  string
	ResourceMeasures string
}

// Configure configures all the values needed to run the acceptance tests.
// This should be called before run.
func Configure(cfg Configuration) error {

	api = cfg.API
	uapi = cfg.UnauthorizedAPI
	product = cfg.Product
	region = cfg.Region
	plan = cfg.Plan
	planFeatures = cfg.PlanFeatures
	newPlan = cfg.NewPlan
	newPlanFeatures = cfg.NewPlanFeatures

	clientID = cfg.ClientID
	clientSecret = cfg.ClientSecret
	connectorPort = cfg.Port

	var err error
	if cfg.CallbackTimeout != "" {
		cbTimeout, err = time.ParseDuration(cfg.CallbackTimeout)
		if err != nil {
			return err
		}

		if cbTimeout > maxTimeout {
			return errors.New("callback timeout cannot exceed 24hrs")
		}
	}

	if cfg.ResourceMeasures != "" {
		err := json.Unmarshal([]byte(cfg.ResourceMeasures), &resourceMeasures)
		if err != nil {
			return errors.Wrap(err, "failed to parse resource measures json")
		}
	}

	fakeConnector, err = connector.New(connectorPort, clientID, clientSecret, product)
	return err
}

var shouldRunErrorCases = true

// Run runs the acceptance tests for all features, less the ones with labels
// in exclude.
//
// run returns a bool indicating success or failure
func Run(ctx context.Context, runErrorCases bool, exclude []string) bool {
	shouldRunErrorCases = runErrorCases
	gomega.RegisterFailHandler(failHandler)
	fakeConnector.Start()
	defer fakeConnector.Stop()

	ok := walkGraph(ctx, exclude, false, execute)
	printSummary(failures, success)
	return ok
}

// Validate checks if the test run has all the required information it needs to run
// the tests.
//
// Valid returns a slice of unique errors where a feature is missing a required
// testing parameter or setting.
func Validate(ctx context.Context, cctx *cli.Context, exclude []string) []error {
	validationErrors := map[string]error{}
	visitorFunc := func(ctx context.Context, feature *FeatureImpl) bool {
		for _, flag := range feature.requiredFlags {
			if !cctx.IsSet(flag) {
				key := fmt.Sprintf("%s-%s-%s", feature.label, feature.name, flag)
				validationErrors[key] = fmt.Errorf("feature `%s` requires flag `%s` to be set", feature.label, flag)
			}
		}

		// always return true so we walk the full graph and get all the errors
		return true
	}

	walkGraph(ctx, exclude, true, visitorFunc)

	var i int
	errs := make([]error, len(validationErrors))
	for _, err := range validationErrors {
		errs[i] = err
		i++
	}
	return errs
}

// walkGraph builds a graph of the features that we've set up and goes over all
// the children. If a feature is marked to be excluded, the feature is skipped,
// otherwise, the visitorFunc is performed with the Feature Implementation.
//
// walkGraph returns a boolean indicating if there were any errors running
// the visitorFuncs.
func walkGraph(ctx context.Context, exclude []string, descendOnErr bool, visitor visitorFunc) bool {
	root := buildGraph(features)
	stack := root.children
	failures := false

	for len(stack) != 0 {
		var n *node
		n, stack = stack[0], stack[1:]

		if !shouldRun(exclude, n.f.label) {
			continue
		}

		ok := visitor(ctx, n.f)
		if !ok {
			n.f.failed = true
		}

		// we only run children if the parent passed when descendOnErr is false.
		// this handles the RunsInside logic.
		if ok || descendOnErr {
			stack = append(n.children, stack...)
		}

		failures = failures || ok
	}

	return failures
}

const testFailurePanic = "test failed"

func execute(ctx context.Context, f *FeatureImpl) (ok bool) {
	defer func() {
		if r := recover(); r == testFailurePanic {
			ok = false
		} else if r != nil {
			panic(r)
		}
	}()

	ok = true

	if !f.hasRun {
		enter(bold(f.label+": ") + f.name)

		f.hasRun = true
		f.fn(ctx)
	} else {
		exit()
		if f.teardown != nil && !f.failed {
			enter(bold(f.label+": ") + f.teardown.name)

			f.teardown.fn(ctx)
			exit()
		}
	}

	return
}

func shouldRun(excluded []string, label string) bool {
	for _, e := range excluded {
		if label == e {
			return false
		}
	}
	return true
}

func failHandler(message string, callerskip ...int) {
	if message[len(message)-1] != '\n' {
		message = message + "\n"
	}

	printIndented(message)

	// bounce out of executing the rest of the test flow
	panic(testFailurePanic)
}

var errFatal = errors.New("Fatal error")

// FatalErr will cause a feature test to fail, and abort running the rest of it.
func FatalErr(format string, args ...interface{}) error {
	failHandler(fmt.Sprintf(format, args...))
	return errFatal // never actually returns this value, as we panic instead.
}
