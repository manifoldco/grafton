package acceptance

import (
	"context"
	"time"

	"github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"
	gm "github.com/onsi/gomega"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/connector"
)

var rotatedCredentialID manifold.ID

var rotateCreds = Feature("credentials_rotation", "Credential rotation", func(ctx context.Context) {
	// TODO: Change this to a bunch of ifs using the rotation type flag
	rotateCredsMultipleManual(ctx)
})

func rotateCredsMultipleManual(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		cID, creds, callbackID, async, err := provisionCredentials(ctx, api, resourceID)
		gm.Expect(err).To(notError(), "Expected a successful provision of a new set of Credentials")

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
				gm.BeNumerically(">", 0),
				"One or more credential should be returned during provision of a new Credential set",
			)
		}

		gm.Expect(len(creds)).To(
			gm.BeNumerically(">", 0),
			"One or more credentials should be returned during provision of a new Credential set",
		)

		for name := range creds {
			gm.Expect(grafton.ValidCredentialName(name)).To(
				gm.BeTrue(), "Credential name must be of the form "+grafton.NameRegexpString)
		}

		rotatedCredentialID = cID
	})

	ErrorCase("with an invalid resource ID", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error

		fakeResourceID, _ := manifold.NewID(idtype.Resource)
		_, _, _, async, err := provisionCredentials(ctx, api, fakeResourceID)

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

	ErrorCase("with already provisioned credentials - same content acts as created", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error

		_, _, _, async, err := provisionCredentialsID(ctx, api, rotatedCredentialID, resourceID)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Same content should be evaluated during the initial call from Manifold",
		)
		gm.Expect(err).To(
			notError(),
			"Create response should be returned (Repeatable Action)",
		)
	})

	ErrorCase("with a bad signature", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		_, _, _, async, err := provisionCredentials(ctx, uapi, resourceID)

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
}

var _ = rotateCreds.RunsInside("provision")
