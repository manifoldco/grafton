package acceptance

import (
	"context"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
)

var rotationTearDown func(context.Context)

var rotateCreds = Feature("credential-rotation", "Rotate a credential set", func(ctx context.Context) {
	switch credentialType {
	case "single":
		featureReplaceRotation(ctx)
	case "multiple":
		featureSwapRotation(ctx)
	default:
		Default(func() {
			FatalErr("unknown credentialType %s", credentialType)
		})
	}
})

var _ = rotateCreds.TearDown("Remove rotated credential sets", func(ctx context.Context) {
	if rotationTearDown == nil {
		return
	}
	rotationTearDown(ctx)
})

var _ = rotateCreds.RunsInside("provision")

func featureReplaceRotation(ctx context.Context) {
	var rotatedCredentialID manifold.ID
	Case("single credential replace", func() {
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

func featureSwapRotation(ctx context.Context) {
	var rotatedCredentialID manifold.ID
	Case("multiple credentials swap", func() {
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
