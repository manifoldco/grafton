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

		pullResourceMeasures(ctx, api, resourceID, resourceMeasures)
	})
})

var _ = measures.RunsInside("provision")

func pullResourceMeasures(ctx context.Context, api *grafton.Client,
	rid manifold.ID, measures map[string]int64) {

	year, month, _ := time.Now().UTC().Date()
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	rm, err := api.PullResourceMeasures(ctx, rid, start, end)

	gm.Expect(err).To(notError(), "No error is expected")

	gm.Expect(rm.ResourceID).To(gm.Equal(rid))

	gm.Expect(rm.PeriodStart).ToNot(gm.BeNil())
	gm.Expect(time.Time(*rm.PeriodStart)).To(gm.Equal(start))

	gm.Expect(rm.PeriodEnd).ToNot(gm.BeNil())
	gm.Expect(time.Time(*rm.PeriodEnd)).To(gm.Equal(end))

	gm.Expect(rm.Measures).To(gm.Equal(resourceMeasures))
}
