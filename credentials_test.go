package grafton

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestValidCredentialName(t *testing.T) {
	tcs := []struct {
		in  string
		out bool
	}{
		{"090asdfaf", false},
		{"AS_DSDA", true},
		{"_DDF", false},
		{"DFS_SDF", true},
		{"D08_DF_", true},
	}

	for _, tc := range tcs {
		t.Run(tc.in, func(t *testing.T) {
			gm.RegisterTestingT(t)

			gm.Expect(ValidCredentialName(tc.in)).To(gm.Equal(tc.out))
		})
	}
}
