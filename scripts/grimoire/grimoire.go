// Package grimoire provides standard magefile targets.
package grimoire

import (
	"github.com/magefile/mage/mg"

	"github.com/manifoldco/logo/grimoire/cast"
)

// Some convenience methods
var (
	Go   = cast.DirPrepSh("go")
	test = cast.DirPrepSh("CGO_ENABLED=0 go test")
	run  = cast.DirPrepSh("go run")
)

var modules = []string{""} // empty string will be the current directory.

// Vendor installs dependencies into the vendor dir. While we typically use go
// modules for building, sometimes you want to install deps for legacy tooling
// to detect.
func Vendor() error {
	for _, m := range modules {
		if err := Go(m, "mod vendor"); err != nil {
			return err
		}
	}
	return nil
}

// Generate runs go generate to generate any generated files. We check in our
// generated files, so this command only has to be run when you change
// something.
func Generate() error {
	for _, m := range modules {
		if err := Go(m, "generate -v ./..."); err != nil {
			return err
		}
	}
	return nil
}

// Tidy runs tidy on all go modules.
func Tidy() error {
	for _, m := range modules {
		if err := Go(m, "mod tidy"); err != nil {
			return err
		}
	}
	return nil
}

// Test runs our unit tests.
func Test() error {
	for _, m := range modules {
		if err := test(m, "-v ./..."); err != nil {
			return err
		}
	}
	return nil
}

// Lint runs golangci-lint to lint our code.
func Lint() error {
	for _, m := range modules {
		if err := run(m, "github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m ./..."); err != nil {
			return err
		}
	}
	return nil
}

// Cover runs tests with coverage reporting
// depending on your codecov reporting setup, this will just work.
func Cover() error {
	for _, m := range modules {
		if err := test(m, "-covermode=atomic -coverprofile=cover.out ./..."); err != nil {
			return err
		}
	}
	return nil
}

// Ci runs all tasks for continuous integration
func Ci() { mg.SerialDeps(Lint, Cover) }

// ShortCi runs all short tests for continuous integration
func ShortCi() { mg.SerialDeps(Lint, ShortCover) }

// ShortCover runs Cover but only for short tests
func ShortCover() error {
	for _, m := range modules {
		if err := test(m, "-short -covermode=atomic -coverprofile=cover.out ./..."); err != nil {
			return err
		}
	}
	return nil
}

// WithSubmodule tells our mage setup that the calling repo has a submodule.
// It's probably a client module. If this is called, linting, testing, and
// generation will be called on the main module and the submodule.
//
// This function can be called multiple times.
func WithSubmodule(dir string) struct{} {
	modules = append(modules, dir)
	return struct{}{}
}

// Modules returns all registered module paths.
func Modules() []string {
	return modules
}
