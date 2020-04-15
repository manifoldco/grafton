package connector

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoo/bone"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/errors"

	"github.com/manifoldco/grafton"
)

var (
	errInvalidContentType = grafton.NewError(errors.BadRequestError, "Invalid Content-Type; expected application/json")
	errInvalidGrant       = grafton.NewError(errors.UnauthorizedError, "Invalid Grant")
	errInvalidCBID        = grafton.NewError(errors.BadRequestError, "Invalid Callback ID Provided")
	errCBNotFound         = grafton.NewError(errors.NotFoundError, "Callback not found")
	errCBResolved         = grafton.NewError(errors.ConflictError, "Callback already complete")
	errBadReqBody         = grafton.NewError(errors.BadRequestError, "Could not parse request")
)

func processCallbackHandler(c *FakeConnector, capturer *RequestCapturer) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token, err := authorizeRequest(c, r)
		if err != nil {
			respondWithError(rw, err)
			return
		}

		if r.Header.Get("content-type") != "application/json" {
			respondWithError(rw, errInvalidContentType)
			return
		}

		// Only access tokens granted w/ Client Credentials have sufficient
		// scope to answer a callback
		if token.GrantType != ClientCredentialsGrantType {
			respondWithError(rw, errInvalidGrant)
			return
		}

		ID, err := manifold.DecodeIDFromString(bone.GetValue(r, "id"))
		if err != nil {
			respondWithError(rw, errInvalidCBID)
			return
		}

		cb := c.GetCallback(ID)
		if cb == nil {
			respondWithError(rw, errCBNotFound)
			return
		}

		cbReq := &CallbackRequest{}
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(cbReq)
		if err != nil {
			respondWithError(rw, errBadReqBody)
			return
		}

		capturer.capture(cbReq)
		err = c.TriggerCallback(ID, cbReq.State, cbReq.Message, cbReq.Credentials)
		if err == nil {
			respondWithJSON(rw, nil, 204)
			return
		}

		switch err {
		case ErrCallbackNotFound:
			respondWithError(rw, errCBNotFound)
			return
		case ErrCallbackAlreadyResolved:
			respondWithError(rw, errCBResolved)
			return
		default:
			respondWithError(rw, errISE)
		}
	})
}
