package acceptance

import (
	"context"
	"fmt"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/connector"
)

var resize = Feature("plan-change", "Change a resource's plan", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		attemptResize(ctx, api, resourceID, newPlan)
	})

	ErrorCase("with existing plan - returns success", func() {
		_, async, err := changePlan(ctx, api, resourceID, newPlan)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Same content should be evaluated during the initial call from Manifold",
		)
		gm.Expect(err).To(notError(),
			"If the current plan matches the requested plan a 204 No Content should be returned (Repeatable Action)")
	})

	ErrorCase("with a non existing resource", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		fakeID, _ := manifold.NewID(idtype.Resource)
		_, async, err := changePlan(ctx, api, fakeID, newPlan)

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
		gm.Expect(e.Type).Should(gm.Equal(merrors.NotFoundError))
	})

	ErrorCase("with a non existing plan", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		_, async, err := changePlan(ctx, api, resourceID, "non-existing")

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

		_, async, err := changePlan(ctx, uapi, resourceID, newPlan)

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
})

var _ = resize.RunsInside("provision")
var _ = resize.RunsBefore("credentials")
var _ = resize.RequiredFlags("new-plan")

var _ = resize.TearDown("Change the resource's plan back to the original", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		attemptResize(ctx, api, resourceID, plan)
	})
})

func attemptResize(ctx context.Context, api *grafton.Client, resourceID manifold.ID, newPlan string) {
	callbackID, async, err := changePlan(ctx, api, resourceID, newPlan)

	gm.Expect(err).To(notError(), "Expected a successful plan change of a resource")

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
			"Credentials cannot be returned on a resource plan change callback",
		)
	}
}

func changePlan(ctx context.Context, api *grafton.Client, resourceID manifold.ID, newPlan string) (manifold.ID, bool, error) {
	Infof("Attempting to resize resource %s to %s\n", resourceID, newPlan)

	c, err := fakeConnector.AddCallback(connector.ResourceResizeCallback)
	if err != nil {
		return manifold.ID{}, false, err
	}

	msg, callback, err := api.ChangePlan(ctx, c.ID, resourceID, newPlan)
	if err != nil {
		return c.ID, callback, err
	}

	if callback {
		Infoln(fmt.Sprintf("Waiting for callback (max: %.1f minutes): %s", cbTimeout.Minutes(), msg))

		cb, err := waitForCallback(c.ID, cbTimeout)
		if err != nil {
			return c.ID, callback, err
		}

		msg = cb.Message
	}

	Infoln("Successfully resized!")
	if msg != "" {
		Infoln("Message: ", msg)
	}

	return c.ID, callback, nil
}
