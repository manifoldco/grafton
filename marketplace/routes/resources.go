package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-zoo/bone"
	"github.com/manifoldco/go-manifold/idtype"
	"github.com/manifoldco/grafton"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/names"

	"github.com/manifoldco/grafton/connector"
	"github.com/manifoldco/grafton/db"
)

func respondResourcePage(d *db.DB, rw http.ResponseWriter, req *http.Request, code int, features string) {
	// Get all resources
	rs := make([]db.Resource, len(d.ResourcesByID))
	i := 0
	for _, r := range d.ResourcesByID {
		rs[i] = r
		i++
	}

	content := struct {
		Resources []db.Resource
		Code      int
		Features  string
	}{
		Resources: rs,
		Code:      code,
		Features:  features,
	}

	respond(rw, req, "resources", content, code)
}

// GetResourcesHandler displays a list of resources in format depending on the
//  Accept header
func GetResourcesHandler(d *db.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		fjson := query.Get("features")

		respondResourcePage(d, rw, req, 200, fjson)
	}
}

// PostResourcesHandler attempts to provision a new resource
func PostResourcesHandler(d *db.DB, gc *grafton.Client,
	fc *connector.FakeConnector) http.HandlerFunc {

	return func(rw http.ResponseWriter, req *http.Request) {
		id, err := manifold.NewID(idtype.Resource)
		if err != nil {
			respondError(rw, req, "Failed to generate ID for resource - "+err.Error(), 500)
			return
		}

		err = req.ParseForm()
		if err != nil {
			respondError(rw, req, "Failed to parse form - "+err.Error(), 500)
			return
		}

		var features manifold.FeatureMap

		featuresTxt := req.Form.Get("features")
		if featuresTxt != "" {
			err = json.Unmarshal([]byte(featuresTxt), &features)

			if err != nil {
				respondError(rw, req, "Failed to parse features - "+err.Error(), 500)
				return
			}
		}

		name := names.ForResource(manifold.Label("grafton"), id)

		// Store in a provisioning state
		r := &db.Resource{
			ID:    id,
			Name:  manifold.Name(name),
			Label: manifold.Label(name),
			// TODO: use real values from config
			Plan:      "ursa-minor",
			Product:   "bear",
			Region:    "all::global",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			State:     db.ResourceStateProvisioning,
			Features:  features,
		}
		d.PutResource(*r)

		failedToProvision := func(m string) {
			r.State = db.ResourceStateProvisionFailed
			d.PutResource(*r)
			respondError(rw, req, m, 500)
		}

		// Request to provision
		cb, err := fc.AddCallback(connector.ResourceProvisionCallback)
		if err != nil {
			failedToProvision("Failed to register callback for resource provision - " + err.Error())
			return
		}

		_, callback, err := gc.ProvisionResource(req.Context(), cb.ID, grafton.ResourceBody{
			ID:       r.ID,
			Product:  string(r.Product),
			Plan:     string(r.Plan),
			Region:   string(r.Region),
			Features: r.Features,
		})

		if callback {
			waitForCallback(fc, cb.ID)
			if cb.State != connector.DoneCallbackState {
				failedToProvision("Failed to provision resource from provider - " + cb.Message)
				return
			}
			// else all-good fall-through
		} else if err != nil {
			failedToProvision("Failed to provision resource from provider - " + err.Error())
			return
		}

		// Provisioned!
		r.State = db.ResourceStateProvisioned
		d.PutResource(*r)

		http.Redirect(rw, req, "/", http.StatusFound)
	}
}

// PutResourcesHandler attempts to update an existing resource
func PutResourcesHandler(d *db.DB) http.HandlerFunc {
	return nil
}

// DeleteResourcesHandler attempts to update an existing resource
func DeleteResourcesHandler(d *db.DB, gc *grafton.Client,
	fc *connector.FakeConnector) http.HandlerFunc {

	return func(rw http.ResponseWriter, req *http.Request) {

		idString := bone.GetValue(req, "id")
		if idString == "" {
			respondError(rw, req, "No ID provided!", 400)
			return
		}

		id, err := manifold.DecodeIDFromString(idString)
		if err != nil {
			respondError(rw, req, "Provided ID was not a Manifold ID", 400)
			return
		} else if id.Type() != idtype.Resource {
			respondError(rw, req, "Provided ID is not for Resource", 400)
			return
		}

		r := d.GetResource(id)
		if r == nil {
			respondError(rw, req, "Resource does not exist", 404)
			return
		}

		originalState := r.State
		failedToDeprovision := func(m string) {
			// Restore state and error
			r.State = originalState
			d.PutResource(*r)
			respondError(rw, req, m, 500)
		}

		// Set resource as deprovisioning
		r.State = db.ResourceStateDerovisioning
		d.PutResource(*r)

		// Request to deprovision
		cb, err := fc.AddCallback(connector.ResourceDeprovisionCallback)
		if err != nil {
			failedToDeprovision("Failed to register callback for resource deprovision - " + err.Error())
			return
		}

		_, callback, err := gc.DeprovisionResource(req.Context(), cb.ID, id)

		if callback {
			waitForCallback(fc, cb.ID)
			if cb.State != connector.DoneCallbackState {
				failedToDeprovision("Failed to deprovision resource from provider - " + cb.Message)
				return
			}
			// else all-good fall-through
		} else if err != nil {
			failedToDeprovision("Failed to deprovision resource from provider - " + err.Error())
			return
		}

		// Deprovisioned!
		r.State = db.ResourceStateDeprovisioned
		d.PutResource(*r)

		http.Redirect(rw, req, "/", http.StatusFound)
	}
}

// SSOResourcesHandler attempts to update an existing resource
func SSOResourcesHandler(d *db.DB, gc *grafton.Client,
	fc *connector.FakeConnector) http.HandlerFunc {

	return func(rw http.ResponseWriter, req *http.Request) {

		idString := bone.GetValue(req, "id")
		if idString == "" {
			respondError(rw, req, "No ID provided!", 400)
			return
		}

		id, err := manifold.DecodeIDFromString(idString)
		if err != nil {
			respondError(rw, req, "Provided ID was not a Manifold ID", 400)
			return
		} else if id.Type() != idtype.Resource {
			respondError(rw, req, "Provided ID is not for Resource", 400)
			return
		}

		authCode, err := fc.CreateCode()
		if err != nil {
			respondError(rw, req, "Failed to create auth code: "+err.Error(), 500)
			return
		}

		url := gc.CreateSsoURL(authCode.Code, id)

		// Redirect to provider dashboard
		rw.Header().Add("Location", url.String())
		rw.WriteHeader(302)
	}
}
