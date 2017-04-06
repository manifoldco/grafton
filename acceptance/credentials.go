package acceptance

import (
	"context"
	"fmt"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/idtype"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/connector"
	"github.com/manifoldco/grafton/generated/provider/client/credential"
)

var credentialID manifold.ID

var creds = Feature("credentials", "Create a credential set", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		cID, creds, callbackID, async, err := provisionCredentials(ctx, api, resourceID)
		gm.Expect(err).To(gm.BeNil(), "Expected a successful provision of Credentials")

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
				"One or more credential should be returned during provision of a Credential set",
			)
		}

		gm.Expect(len(creds)).To(
			gm.BeNumerically(">", 0),
			"One or more credentials should be returned during provision of a Credential set",
		)

		credentialID = cID
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
			gm.BeAssignableToTypeOf(&credential.PutCredentialsIDNotFound{}),
			"Expected an IDNotFound error, got %T", err,
		)
	})

	ErrorCase("with already provisioned credentials - same content acts as created", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		var err error

		_, _, _, async, err := provisionCredentialsID(ctx, api, credentialID, resourceID)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Same content should be evaluated during the initial call from Manifold",
		)
		gm.Expect(err).To(
			gm.BeNil(),
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
			gm.BeAssignableToTypeOf(&credential.PutCredentialsIDUnauthorized{}),
			"Expected an Unauthorized error, got %T", err,
		)
	})
})

var _ = creds.TearDown("Delete a credential set", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		callbackID, async, err := deprovisionCredentials(ctx, api, credentialID)
		gm.Expect(err).To(gm.BeNil(), "No error is expected")

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
				"Credentials cannot be returned on a deprovisioning callback",
			)
		}
	})

	ErrorCase("delete an unexisting credential", func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		fakeCredentialID, _ := manifold.NewID(idtype.Credential)
		_, async, err := deprovisionCredentials(ctx, api, fakeCredentialID)

		gm.Expect(async).To(
			gm.BeFalse(),
			"Validation errors should be returned on the initial request",
		)
		gm.Expect(err).ShouldNot(
			gm.BeNil(),
			"Expected an error, got nil",
		)
		gm.Expect(err).Should(
			gm.BeAssignableToTypeOf(&credential.DeleteCredentialsIDNotFound{}),
			"Expected a CredentialsIDNotFound error, got %T", err,
		)
	})
})

var _ = creds.RunsInside("provision")

func provisionCredentials(ctx context.Context, api *grafton.Client, resourceID manifold.ID) (manifold.ID, map[string]string, manifold.ID, bool, error) {
	Infof("Attempting to provision credentials for resource: %s\n", resourceID)
	ID, err := manifold.NewID(idtype.Credential)
	if err != nil {
		return ID, nil, ID, false, FatalErr("Could not generate credential id:", err)
	}

	return provisionCredentialsID(ctx, api, ID, resourceID)
}

func provisionCredentialsID(ctx context.Context, api *grafton.Client, credentialID, resourceID manifold.ID) (manifold.ID, map[string]string, manifold.ID, bool, error) {
	c, err := fakeConnector.AddCallback(connector.CredentialProvisionCallback)
	if err != nil {
		return credentialID, nil, manifold.ID{}, false, err
	}

	creds, msg, callback, err := api.ProvisionCredentials(ctx, c.ID, resourceID, credentialID)
	if err != nil {
		return credentialID, nil, c.ID, false, err
	}

	if callback {
		Infoln(fmt.Sprintf("Waiting for Callback: (max: %.1f minutes): %s", cbTimeout.Minutes(), msg))

		cb, err := waitForCallback(c.ID, cbTimeout)
		if err != nil {
			return credentialID, nil, c.ID, callback, err
		}

		msg = cb.Message
		creds = cb.Credentials
	}

	Infoln("Provisioned Credentials Successfully")
	if msg != "" {
		Infoln("Message: ", msg)
	}
	Infoln("Credentials:")
	for k, v := range creds {
		Infoln("  ", k, "=", v)
	}

	return credentialID, creds, c.ID, callback, err
}

func deprovisionCredentials(ctx context.Context, api *grafton.Client, credentialID manifold.ID) (manifold.ID, bool, error) {
	Infoln("Attempting to deprovision credentials:", credentialID)

	c, err := fakeConnector.AddCallback(connector.CredentialDeprovisionCallback)
	if err != nil {
		return manifold.ID{}, false, err
	}

	msg, callback, err := api.DeprovisionCredentials(ctx, c.ID, credentialID)
	if err != nil {
		return c.ID, callback, err
	}

	if callback {
		Infoln(fmt.Sprintf("Waiting for Callback(max %.1f minutes): %s", cbTimeout.Minutes(), msg))

		cb, err := waitForCallback(c.ID, cbTimeout)
		if err != nil {
			return c.ID, callback, err
		}

		msg = cb.Message
	}

	Infoln("Credential Deprovisioned.")
	if msg != "" {
		Infoln("Message: ", msg)
	}
	return c.ID, callback, nil
}
