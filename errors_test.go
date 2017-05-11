package grafton

import (
	"errors"
	"net/http"
	"testing"

	swagerrs "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"
)

func TestIsFatal(t *testing.T) {
	t.Run("a generic error is not fatal", func(t *testing.T) {
		gm.RegisterTestingT(t)

		err := errors.New("hi")
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
	})

	t.Run("an InternalServerError is not fatal", func(t *testing.T) {
		gm.RegisterTestingT(t)

		err := manifold.NewError(merrors.InternalServerError, "hi")
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
	})

	t.Run("a go-openapi with code < 500 is fatal", func(t *testing.T) {
		gm.RegisterTestingT(t)

		err := swagerrs.New(401, "hi")
		gm.Expect(IsFatal(err)).To(gm.BeTrue())
	})

	t.Run("a go-openapi with code 5xx is not fatal", func(t *testing.T) {
		gm.RegisterTestingT(t)

		err := swagerrs.New(503, "hi")
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
	})

	t.Run("a go-openapi with code 4xx is fatal", func(t *testing.T) {
		gm.RegisterTestingT(t)

		err := swagerrs.New(401, "what")
		gm.Expect(IsFatal(err)).To(gm.BeTrue())
	})

	t.Run("a runtime APIError with code 5xx is not fatal", func(t *testing.T) {
		gm.RegisterTestingT(t)

		err := runtime.NewAPIError("unknown", nil, 500)
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
	})

	t.Run("a runtime APIError with code 4xx to be fatal", func(t *testing.T) {
		gm.RegisterTestingT(t)

		err := runtime.NewAPIError("unknown", nil, 400)
		gm.Expect(IsFatal(err)).To(gm.BeTrue())
	})
}

func TestToError(t *testing.T) {
	t.Run("with an Error", func(t *testing.T) {
		derr := NewError(merrors.BadRequestError, "Bad request")

		if rerr := ToError(derr); rerr != derr {
			t.Errorf("Expected %#v to equal %#v", derr, rerr)
		}
	})

	t.Run("with an HTTPError", func(t *testing.T) {
		t.Run("with a known error code", func(t *testing.T) {
			err := &mockHTTPError{400, "Test error"}
			derr := ToError(err)
			rerr, ok := derr.(*Error)
			if !ok {
				t.Errorf("Expected %#v to be of type `*Error`; it's not", derr)
			}

			if et, _ := merrors.TypeForStatusCode(err.code); et != rerr.Type {
				t.Errorf("Expected Type to equal `%s`, got `%s`", et, rerr.Type)
			}

			if rerr.Message != err.message {
				t.Errorf("Expected Message to equal `%s`, got `%s`", err.message, rerr.Message)
			}
		})

		t.Run("with an unknown error code", func(t *testing.T) {
			err := &mockHTTPError{200, "Test error"}
			derr := ToError(err)
			rerr, ok := derr.(*Error)
			if !ok {
				t.Errorf("Expected %#v to be of type `*Error`; it's not", derr)
			}

			if rerr.Type != merrors.InternalServerError {
				t.Errorf("Expected Typeto be `%s`, got `%s`", merrors.InternalServerError, rerr.Type)
			}
		})
	})

	t.Run("with an unknown error type", func(t *testing.T) {
		err := ToError(errors.New("test"))
		derr, ok := err.(*Error)

		if !ok {
			t.Errorf("Expected %#v to be of type `merrors.ProviderError`, it's not", err)
		}

		if derr.Type != merrors.InternalServerError {
			t.Errorf("Expected ErrorType to be `%s`, got `%s`", merrors.InternalServerError, derr.Type)
		}
	})
}

type mockHTTPError struct {
	code    int
	message string
}

func (e *mockHTTPError) Error() string {
	return e.message
}

func (e *mockHTTPError) StatusCode() int {
	return e.code
}

func (e *mockHTTPError) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	pr.Produce(rw, e)
}
