package acceptance

import (
	"context"
	"time"
)

var _ = Feature("cleanup", "Remove dangling resource due to failed provision", func(ctx context.Context) {
	Default(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		curResource := attemptResourceProvision(ctx, api, product, plan, region)
		attemptResourceDeprovision(ctx, api, curResource.ID)
	})
})
