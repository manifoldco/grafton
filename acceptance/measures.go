package acceptance

import (
	"context"
	"time"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/grafton"
	gm "github.com/onsi/gomega"
)

var measures = Feature("resource-measures", "Pull usage measures from a Resource", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		pullResourceMeasures(ctx, api, resourceID)
	})
})

var _ = measures.RunsInside("provision")

func pullResourceMeasures(ctx context.Context, api *grafton.Client, rid manifold.ID) {
	year, month, _ := time.Now().UTC().Date()
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	rm, err := api.PullResourceMeasures(ctx, rid, start, end)

	gm.Expect(err).To(notError(), "No error is expected")

	gm.Expect(rm.PeriodStart).ToNot(gm.BeNil())
}
