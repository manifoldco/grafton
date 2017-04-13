package acceptance

import (
	"context"
	"testing"
)

func TestRequiredFlags(t *testing.T) {
	t.Run("it should add a flag", func(t *testing.T) {
		feature := Feature("test", "Required Flags Test", func(context.Context) {})
		feature.RequiredFlags("my-test")

		if !feature.NeedsFlag("my-test") {
			t.Errorf("Feature should need `my-test` flag")
		}
	})

	t.Run("it should add multiple flags", func(t *testing.T) {
		feature := Feature("test", "Required Flags Test", func(context.Context) {})
		feature.RequiredFlags("my-test", "my-second-test")

		if !feature.NeedsFlag("my-test") {
			t.Errorf("Feature should need `my-test` flag")
		}

		if !feature.NeedsFlag("my-second-test") {
			t.Errorf("Feature should need `my-second-test` flag")
		}
	})
}
