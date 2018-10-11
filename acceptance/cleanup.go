package acceptance

import (
	"context"
	"time"
)

var _ = Feature("cleanup", "Can provision and deprovision a resource", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		curResource := attemptResourceProvision(ctx, api, product, plan, planFeatures, region)
		attemptResourceDeprovision(ctx, api, curResource.ID)
	})
})
