package grafton

import (
	"errors"
	"net/http"

	swagerrs "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"
)

// ErrMissingMsg occurs when a provider's response is missing the Message field
var ErrMissingMsg = errors.New("`message` field was missing from the response")

// NewErrWithMsg creates a new Error from a string pointer, if the pointer is
// nil then an ErrMissingMsg is returned instead.
func NewErrWithMsg(t merrors.Type, m *string) error {
	if m == nil {
		return ErrMissingMsg
	}

	return NewError(t, *m)
}

// IsFatal returns true or false depending on whether or not the given error is
// considered to be fatal
func IsFatal(err error) bool {
	switch e := err.(type) {
	case *Error:
		if e.Type == merrors.InternalServerError {
			return false
		}

		return true
	case swagerrs.Error:
		// If status code is 5xx, then something went terribly wrong on their
		// side, so we want to try again
		//
		// 6xx and above are status codes used by go-openapi to represent
		// different schema related errors
		//
		// We don't consider 2xx codes as non-fatal because someplace upstream
		// we didn't understand what was sent to us.
		code := int(e.Code())
		if code >= 500 && code < 600 {
			return false
		}

		return true
	case *runtime.APIError:
		// If its a runtime error (which occurs inside the client) has a status
		// code between 500 and 600, then it's not fatal, otherwise, it is :)
		//
		// We don't consider 2xx codes as non-fatal because someplace upstream
		// we didn't understand what was sent to us.
		if e.Code >= 500 && e.Code < 600 {
			return false
		}

		return true
	default:
		return false
	}
}

// Error is the error type we'll use for providers to return.
type Error struct {
	Type    merrors.Type `json:"error"`
	Message string       `json:"message"`
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
		if et, ok := merrors.TypeForStatusCode(e.StatusCode()); ok {
			return NewError(et, e.Error())
		}
	}

	return NewError(merrors.InternalServerError, "Internal Server Error")
}

// NewError creates a new Error with a type and message.
func NewError(t merrors.Type, m string) *Error {
	return &Error{
		Type:    t,
		Message: m,
	}
}
