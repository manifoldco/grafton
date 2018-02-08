package acceptance

import (
	"context"
	"errors"
	"time"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/grafton"
	gm "github.com/onsi/gomega"
)

var measures = Feature("resource-measures", "Pull usage measures from a Resource", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		err := pullResourceMeasures(ctx, api, resourceID)
		gm.Expect(err).To(notError(), "No error is expected")
	})
})

var _ = measures.RunsInside("provision")

func pullResourceMeasures(ctx context.Context, api *grafton.Client, rid manifold.ID) error {
	return errors.New("not implemented")
}
