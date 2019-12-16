package acceptance

import (
	"context"
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()

	t.Run("without any flags set", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		cctx := cli.NewContext(nil, set, nil)

		t.Run("with all features enabled", func(t *testing.T) {
			errs := Validate(ctx, cctx, []string{})

			// get the number of flags for all features
			totalRequiredFlags := 0
			for _, f := range features {
				totalRequiredFlags += len(f.requiredFlags)
			}

			if len(errs) != totalRequiredFlags {
				t.Errorf("Expected `%d` errors, got `%d`", totalRequiredFlags, len(errs))
			}
		})

		t.Run("with a feature excluded", func(t *testing.T) {
			errs := Validate(ctx, cctx, []string{"sso"})

			// get the number of flags for all features
			totalRequiredFlags := 0
			for _, f := range features {
				if f.label != "sso" {
					totalRequiredFlags += len(f.requiredFlags)
				}
			}

			if len(errs) != totalRequiredFlags {
				t.Errorf("Expected `%d` errors, got `%d`", totalRequiredFlags, len(errs))
			}
		})
	})

	t.Run("with flags for a feature set", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		// required flags for the plan-change test
		set.Set("new-plan", "my-new-plan")
		cctx := cli.NewContext(nil, set, nil)

		t.Run("with a feature excluded", func(t *testing.T) {
			errs := Validate(ctx, cctx, []string{"plan-change"})

			// get the number of flags for all features
			totalRequiredFlags := 0
			for _, f := range features {
				if f.label != "plan-change" {
					totalRequiredFlags += len(f.requiredFlags)
				}
			}

			if len(errs) != totalRequiredFlags {
				t.Errorf("Expected `%d` errors, got `%d`", totalRequiredFlags, len(errs))
			}
		})
	})
}
