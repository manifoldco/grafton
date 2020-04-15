// +build tools

package tools

import (
	"github.com/client9/misspell/cmd/misspell" // lint
	"github.com/go-swagger/go-swagger/cmd/swagger"
	"github.com/golangci/golangci-lint/cmd/golangci-lint" // lint
	"github.com/gordonklaus/ineffassign"                  // lint
	"github.com/tsenart/deadcode"                         // lint
	"golang.org/x/lint/golint"                            // lint
	"github.com/magefile/mage"
	"github.com/gobuffalo/packr"
)
