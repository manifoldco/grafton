package acceptance

import (
	"context"
	"errors"
	"fmt"
	"time"

	gm "github.com/onsi/gomega"

	manifold "github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"
	"github.com/manifoldco/go-manifold/names"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/connector"
	"github.com/manifoldco/grafton/db"
)

var errTimeout = errors.New("Exceeded Callback Wait time")
var resourceID manifold.ID
var curResource *db.Resource

var provision = Feature("provision", "Provision a resource", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		curResource := attemptResourceProvision(ctx, api, product, plan, planFeatures, region)
		resourceID = curResource.ID
	})

	ErrorCase("with a faulty product name", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error
		_, _, async, err := provisionResource(ctx, api, "not-your-product", plan, planFeatures, region)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Validation errors should be returned on the initial request",
		)
		gm.Expect(err).ShouldNot(
			gm.BeNil(),
			"Expected an error, got nil",
		)
		gm.Expect(err).Should(
			gm.BeAssignableToTypeOf(&grafton.Error{}),
			"Expected a grafton error, got %T", err,
		)

		e := err.(*grafton.Error)
		gm.Expect(e.Type).Should(gm.Equal(merrors.BadRequestError))
	})

	ErrorCase("with a faulty plan name", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error
		_, _, async, err := provisionResource(ctx, api, product, "faulty-plan-name", nil, region)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Validation errors should be returned on the initial request",
		)
		gm.Expect(err).ShouldNot(
			gm.BeNil(),
			"Expected an error, got nil",
		)
		gm.Expect(err).Should(
			gm.BeAssignableToTypeOf(&grafton.Error{}),
			"Expected a grafton error, got %T", err,
		)

		e := err.(*grafton.Error)
		gm.Expect(e.Type).Should(gm.Equal(merrors.BadRequestError))
	})

	ErrorCase("with a faulty region", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error
		_, _, async, err := provisionResource(ctx, api, product, plan, planFeatures, "faulty-region")

		gm.Expect(async).To(
			gm.BeFalse(),
			"Validation errors should be returned on the initial request",
		)
		gm.Expect(err).ShouldNot(
			gm.BeNil(),
			"Expected an error, got nil",
		)
		gm.Expect(err).Should(
			gm.BeAssignableToTypeOf(&grafton.Error{}),
			"Expected a grafton error, got %T", err,
		)

		e := err.(*grafton.Error)
		gm.Expect(e.Type).Should(gm.Equal(merrors.BadRequestError))
	})

	ErrorCase("with a bad signature", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error
		_, _, async, err := provisionResource(ctx, uapi, product, plan, planFeatures, region)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Validation errors should be returned on the initial request",
		)
		gm.Expect(err).ShouldNot(
			gm.BeNil(),
			"Expected an error, got nil",
		)

		gm.Expect(err).Should(
			gm.BeAssignableToTypeOf(&grafton.Error{}),
			"Expected a grafton error, got %T", err,
		)

		e := err.(*grafton.Error)
		gm.Expect(e.Type).Should(gm.Equal(merrors.UnauthorizedError))
	})

	ErrorCase("with an already provisioned resource - same content acts as created", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error
		_, callbackID, async, err := provisionResourceID(ctx, api, resourceID, product, plan, planFeatures, region)

		if async {
			c := fakeConnector.GetCallback(callbackID)

			gm.Expect(c.State).To(
				gm.Equal(connector.DoneCallbackState),
				"Expected to receive 'done' as the state",
			)
		}

		gm.Expect(err).To(notError(), "Create response should be returned (Repeatable Action)")
	})

	ErrorCase("with an already provisioned resource - different content results in conflict", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error
		_, callbackID, async, err := provisionResourceID(ctx, api, resourceID, product, newPlan, newPlanFeatures, region)

		if async {
			c := fakeConnector.GetCallback(callbackID)

			gm.Expect(c.State).To(
				gm.Equal(connector.ErrorCallbackState),
				"Expected to receive 'error' as the state",
			)
		}

		gm.Expect(err).ShouldNot(
			gm.BeNil(),
			"Expected an error, got nil",
		)
		gm.Expect(err).Should(
			gm.BeAssignableToTypeOf(&grafton.Error{}),
			"Expected a grafton error, got %T", err,
		)

		e := err.(*grafton.Error)
		gm.Expect(e.Type).Should(gm.Equal(merrors.ConflictError))
	})
})

var _ = provision.TearDown("Deprovision a resource", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		attemptResourceDeprovision(ctx, api, resourceID)
	})

	ErrorCase("delete a non existing resource", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		fakeID, _ := manifold.NewID(idtype.Resource)
		_, async, err := deprovisionResource(ctx, api, fakeID)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Resource existence should be evaluated during the initial call from Manifold",
		)
		gm.Expect(err).ShouldNot(
			gm.BeNil(),
			"Expected an error, got nil",
		)
		gm.Expect(err).Should(
			gm.BeAssignableToTypeOf(&grafton.Error{}),
			"Expected a grafton error, got %T", err,
		)

		e := err.(*grafton.Error)
		gm.Expect(e.Type).Should(gm.Equal(merrors.NotFoundError))
	})
})
var _ = provision.RequiredFlags("product", "plan", "region", "new-plan")

func attemptResourceProvision(ctx context.Context, api *grafton.Client, product, plan string,
	planFeatures manifold.FeatureMap, region string) *db.Resource {

	var err error
	curResource, callbackID, async, err := provisionResource(ctx, api, product, plan, planFeatures, region)
	gm.Expect(err).To(notError(), "Expected a successful provision of a resource")

	if async {
		c := fakeConnector.GetCallback(callbackID)

		gm.Expect(c.State).To(
			gm.Equal(connector.DoneCallbackState),
			"Expected to receive 'done' as the state",
		)
		gm.Expect(len(c.Message)).To(gm.SatisfyAll(
			gm.BeNumerically(">=", 3),
			gm.BeNumerically("<", 256),
		), "Message must be between 3 and 256 characters long.")
		gm.Expect(len(c.Credentials)).To(
			gm.Equal(0),
			"Credentials cannot be returned on a resource provisioning callback",
		)
	}

	return curResource
}

func attemptResourceDeprovision(ctx context.Context, api *grafton.Client, resourceID manifold.ID) {
	callbackID, async, err := deprovisionResource(ctx, api, resourceID)

	gm.Expect(err).To(notError(), "No error is expected")
	if async {
		c := fakeConnector.GetCallback(callbackID)

		gm.Expect(c.State).To(
			gm.Equal(connector.DoneCallbackState),
			"Expected to receive 'done' as the state",
		)
		gm.Expect(len(c.Message)).To(gm.SatisfyAll(
			gm.BeNumerically(">=", 3),
			gm.BeNumerically("<", 256),
		), "Message must be between 3 and 256 characters long.")
		gm.Expect(len(c.Credentials)).To(
			gm.Equal(0),
			"Credentials cannot be returned on a resource deprovisioning callback",
		)
	}
}

func waitForCallback(ID manifold.ID, max time.Duration) (*connector.Callback, error) {
	timeout := time.After(max)

waitForCallback:
	select {
	case cb := <-fakeConnector.OnCallback:
		if cb.ID != ID {
			goto waitForCallback
		}

		if cb.State == connector.PendingCallbackState {
			goto waitForCallback
		}

		return cb, nil
	case <-timeout:
		return nil, errTimeout
	}
}

func provisionResource(ctx context.Context, api *grafton.Client, product, plan string,
	planFeatures manifold.FeatureMap, region string) (*db.Resource, manifold.ID, bool, error) {

	Infoln("Attempting to provision resource")

	ID, err := manifold.NewID(idtype.Resource)
	if err != nil {
		return nil, ID, false, FatalErr("Could not generate resource id: %s", err)
	}

	return provisionResourceID(ctx, api, ID, product, plan, planFeatures, region)
}

func provisionResourceID(ctx context.Context, api *grafton.Client, id manifold.ID, product, plan string,
	planFeatures manifold.FeatureMap, region string) (*db.Resource, manifold.ID, bool, error) {

	c, err := fakeConnector.AddCallback(connector.ResourceProvisionCallback)
	if err != nil {
		return nil, c.ID, false, err
	}

	productLabel := manifold.Label(product)
	if err := productLabel.Validate(nil); err != nil {
		return nil, c.ID, false, FatalErr("Product label is not a valid label: %s", err)
	}
	planLabel := manifold.Label(plan)
	if err := planLabel.Validate(nil); err != nil {
		return nil, c.ID, false, FatalErr("Plan label is not a valid label: %s", err)
	}

	label := names.ForResource(manifold.Label(product), id)

	r := &db.Resource{
		ID:        id,
		Label:     label,
		Name:      manifold.Name(label),
		Product:   productLabel,
		Plan:      planLabel,
		Region:    region,
		Features:  planFeatures,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Ensure we remove the resource from the connector *if* the resource was
	// not successfully provisioned.
	success := false
	fakeConnector.AddResource(r)
	defer func() {
		if success {
			return
		}

		fakeConnector.RemoveResource(r.ID)
	}()

	model := grafton.ResourceBody{
		ID:       id,
		Product:  product,
		Plan:     plan,
		Region:   region,
		Features: nil,
	}

	msg, callback, err := api.ProvisionResource(ctx, c.ID, model)
	if err != nil {
		return nil, c.ID, false, err
	}

	if callback {
		Infoln(fmt.Sprintf("Waiting for Callback (max: %.1f minutes): %s", cbTimeout.Minutes(), msg))

		cb, err := waitForCallback(c.ID, cbTimeout)
		if err != nil {
			return nil, c.ID, callback, err
		}

		msg = cb.Message
	}

	Infoln("Resource Provisioned Successfully:", id)
	if msg != "" {
		Infoln("Message: ", msg)
	}

	success = true
	return r, c.ID, callback, nil
}

func deprovisionResource(ctx context.Context, api *grafton.Client, resourceID manifold.ID) (manifold.ID, bool, error) {
	Infoln("Attempting to deprovision resource:", resourceID)

	c, err := fakeConnector.AddCallback(connector.ResourceDeprovisionCallback)
	if err != nil {
		return manifold.ID{}, false, err
	}

	msg, callback, err := api.DeprovisionResource(ctx, c.ID, resourceID)
	if err != nil {
		return c.ID, callback, err
	}

	if callback {
		Infoln(fmt.Sprintf("Waiting for Callback (max: %.1f minutes): %s", cbTimeout.Minutes(), msg))

		cb, err := waitForCallback(c.ID, cbTimeout)
		if err != nil {
			return c.ID, callback, err
		}

		msg = cb.Message
	}

	Infoln("Resource Deprovisioned.")
	if msg != "" {
		Infoln("Callback Message: ", msg)
	}

	return c.ID, callback, nil
}
