package acceptance

import (
	"context"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
)

var rotationTearDown func(context.Context)

var rotateCreds = Feature("credentials-rotation", "Credential rotation", func(ctx context.Context) {
	switch credentialRotationType {
	case "manual":
		switch credentialType {
		case "single":
			featureRotateCredsSingleManual(ctx)
		case "multiple":
			featureRotateCredsMultipleManual(ctx)
		default:
			Default(func() {
				FatalErr("unknown credentialType %s", credentialType)
			})
		}
	case "automatic":
		switch credentialType {
		case "single":
			Default(func() {
				FatalErr("automatic single credential rotation test not implemented")
			})
		case "multiple":
			Default(func() {
				FatalErr("automatic multiple credential rotation test not implemented")
			})
		default:
			Default(func() {
				FatalErr("unknown credentialType %s", credentialType)
			})
		}
	case "none":
		// No tests
		Default(func() {})
	default:
		Default(func() {
			FatalErr("unknown credentialRotationType %s", credentialRotationType)
		})
	}
})

var _ = rotateCreds.TearDown("credentials-rotation", func(ctx context.Context) {
	rotationTearDown(ctx)
})

func featureRotateCredsSingleManual(ctx context.Context) {
	var rotatedCredentialID manifold.ID
	Default(func() {
		initialCredID, initialValues := mustProvisionCredentials(ctx, api, resourceID)

		// delete initial credential before creating new one
		mustDeprovisionCredentials(ctx, api, initialCredID)

		rID, rotatedValues := mustProvisionCredentials(ctx, api, resourceID)
		rotatedCredentialID = rID

		// assert initial and rotated are not the same
		gm.Expect(rotatedValues).ToNot(
			gm.Equal(initialValues), "Different credentials expected for new Credential Set")

	})

	rotationTearDown = func(ctx context.Context) {
		Default(func() {
			mustDeprovisionCredentials(ctx, api, rotatedCredentialID)
		})
	}
}

func featureRotateCredsMultipleManual(ctx context.Context) {
	var rotatedCredentialID manifold.ID
	Default(func() {
		initialCredID, initialValues := mustProvisionCredentials(ctx, api, resourceID)
		rID, rotatedValues := mustProvisionCredentials(ctx, api, resourceID)
		rotatedCredentialID = rID

		// assert initial and rotated are not the same
		gm.Expect(rotatedValues).ToNot(
			gm.Equal(initialValues), "Different credentials expected for new Credential Set")

		// delete initial credential
		mustDeprovisionCredentials(ctx, api, initialCredID)
	})

	rotationTearDown = func(ctx context.Context) {
		Default(func() {
			mustDeprovisionCredentials(ctx, api, rotatedCredentialID)
		})
	}
}

var _ = rotateCreds.RunsInside("provision")
