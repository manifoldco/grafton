package connector

import (
	"encoding/json"
	"net/http"

	"github.com/manifoldco/go-manifold/idtype"
	"github.com/manifoldco/grafton/db"

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

		uid, err := manifold.NewID(idtype.User)
		if err != nil {
			respondWithError(rw, err)
			return
		}

		users := []UserTarget{
			{
				ID:    uid,
				Name:  "Manny Fold",
				Email: "manny@manifold.co",
				Role:  UserTargetRoleOwner,
			},
		}

		respondWithJSON(rw, users, 200)
	}
}

func getResourceCredentialsHandler(c *FakeConnector, _ *RequestCapturer) http.HandlerFunc {
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

		creds := c.DB.GetCredentialsByResource(resource.ID)

		respondWithJSON(rw, creds, 200)
	}
}

func getResourceMeasuresHandler(c *FakeConnector, _ *RequestCapturer) http.HandlerFunc {
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

		ms := c.DB.GetMeasuresByResource(ID)
		measures := make([]*ResourceMeasures, len(ms))
		for i, m := range ms {
			rms := &ResourceMeasures{
				UpdatedAt:   m.UpdatedAt,
				PeriodStart: m.PeriodStart,
				PeriodEnd:   m.PeriodEnd,
			}
			for f, v := range m.Measures {
				// TODO: Expand this ( Features, Values ) based on supplied product config
				rms.Measures = append(rms.Measures, ResourceMeasure{
					Feature: ResourceMeasureFeature{
						Name:  manifold.Name(f),
						Label: manifold.Label(f),
					},
					FeatureValue: ResourceMeasureFeatureValue{
						Name:  manifold.Name(f),
						Label: manifold.FeatureValueLabel(f),
					},
					Usage: v,
				})
			}
			measures[i] = rms
		}

		respondWithJSON(rw, measures, 200)
	}
}

func putResourceMeasuresHandler(c *FakeConnector, _ *RequestCapturer) http.HandlerFunc {
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

		inboundMeasure := &db.Measure{}
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(inboundMeasure)
		if err != nil {
			respondWithError(rw, errBadReqBody)
			return
		}

		inboundMeasure.ResourceID = resource.ID
		c.DB.PutMeasure(*inboundMeasure)

		rw.WriteHeader(204)
	}
}
