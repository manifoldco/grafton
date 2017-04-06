package grafton

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/errors"
)

// Error is the error type we'll use for providers to return.
type Error struct {
	Type    errors.Type `json:"error"`
	Message string      `json:"message"`
}

// Error returns the message of the error.
func (e *Error) Error() string {
	return e.Message
}

// StatusCode returns the StatusCode for the current error type.
func (e *Error) StatusCode() int {
	return e.Type.Code()
}

// WriteResponse completes the interface for a HTTPError; enabling an error to
// be returned as a middleware.Responder from go-openapi/runtime
//
// A panic will occur if the given producer errors.
func (e *Error) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(e.StatusCode())
	if err := pr.Produce(rw, e); err != nil {
		panic(err)
	}
}

// ToError receives an error and mutates it into a grafton.Error based
// on the concrete type of the error.
func ToError(err error) manifold.HTTPError {
	switch e := err.(type) {
	case *Error:
		return e
	case manifold.HTTPError:
		if et, ok := errors.TypeForStatusCode(e.StatusCode()); ok {
			return NewError(et, e.Error())
		}
	}

	return NewError(errors.InternalServerError, "Internal Server Error")
}

// NewError creates a new Error with a type and message.
func NewError(t errors.Type, m string) *Error {
	return &Error{
		Type:    t,
		Message: m,
	}
}
