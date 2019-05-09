package main

import (
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

func Test_apiURLPattern(t *testing.T) {
	gm.RegisterTestingT(t)
	oldManifoldHostname := os.Getenv("MANIFOLD_HOSTNAME")
	defer func() {
		os.Setenv("MANIFOLD_HOSTNAME", oldManifoldHostname)
	}()
	oldManifoldScheme := os.Getenv("MANIFOLD_SCHEME")
	defer func() {
		os.Setenv("MANIFOLD_SCHEME", oldManifoldScheme)
	}()

	gm.Expect(apiURLPattern()).To(gm.Equal("https://api.%s.manifold.co/v1"))

	os.Setenv("MANIFOLD_HOSTNAME", "my-hostname")
	os.Setenv("MANIFOLD_SCHEME", "my-scheme")
	gm.Expect(apiURLPattern()).To(gm.Equal("my-scheme://api.%s.my-hostname/v1"))
}
