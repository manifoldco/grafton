package connector

import (
	"net/http"

	"github.com/go-zoo/bone"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/errors"

	"github.com/manifoldco/grafton"
)

var (
	errMissingID       = grafton.NewError(errors.BadRequestError, "Must provide a valid id")
	errInvalidID       = grafton.NewError(errors.BadRequestError, "Invalid Resource ID Provided")
	errMissingResource = grafton.NewError(errors.NotFoundError, "Resource Not Found")
)

func getResourceHandler(c *FakeConnector, _ *RequestCapturer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		_, err := authorizeRequest(c, r)
		if err != nil {
			respondWithError(rw, err)
			return
		}

		idString := bone.GetValue(r, "id")
		if idString == "" {
			respondWithError(rw, errMissingID)
			return
		}

		ID, err := manifold.DecodeIDFromString(idString)
		if err != nil {
			respondWithError(rw, errInvalidID)
			return
		}

		resource := c.GetResource(ID)
		if resource == nil {
			respondWithError(rw, errMissingResource)
			return
		}

		respondWithJSON(rw, resource, 200)
	}
}

func getResourceUsersHandler(c *FakeConnector, _ *RequestCapturer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		_, err := authorizeRequest(c, r)
		if err != nil {
			respondWithError(rw, err)
			return
		}

		idString := bone.GetValue(r, "id")
		if idString == "" {
			respondWithError(rw, errMissingID)
			return
		}

		ID, err := manifold.DecodeIDFromString(idString)
		if err != nil {
			respondWithError(rw, errInvalidID)
			return
		}

		resource := c.GetResource(ID)
		if resource == nil {
			respondWithError(rw, errMissingResource)
			return
		}

		users := []UserTarget{
			{
				Name:  "joe user",
				Email: "joe@user.com",
			},
		}

		respondWithJSON(rw, users, 200)
	}
}
