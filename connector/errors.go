package connector

import (
	"github.com/manifoldco/go-manifold/errors"

	"github.com/manifoldco/grafton"
)

// Shared errors between endpoints
var (
	errISE = grafton.NewError(errors.InternalServerError, "Internal Server Error")
)
