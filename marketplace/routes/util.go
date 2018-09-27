package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/grafton/connector"
	"github.com/manifoldco/grafton/marketplace/pages"
)

const callbackTimeout = time.Minute * 5

func isJSONRequest(req *http.Request) bool {
	ct := req.Header.Get("Content-Type")
	return strings.Contains(strings.ToLower(ct), "/json")
}

func respondWithJSON(rw http.ResponseWriter, v interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	if v == nil {
		return
	}

	enc := json.NewEncoder(rw)
	err := enc.Encode(v)
	if err != nil {
		panic(err)
	}
}

func respondWithHTML(rw http.ResponseWriter, t *template.Template, v interface{}, code int) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(code)

	if v == nil {
		return
	}

	err := t.Execute(rw, v)
	if err != nil {
		panic(err)
	}
}

func respond(rw http.ResponseWriter, req *http.Request, t *template.Template, v interface{}, code int) {
	// If JSON Request, respond as such
	if isJSONRequest(req) {
		respondWithJSON(rw, v, code)
		return
	}
	// Else HTML
	respondWithHTML(rw, t, v, code)
}

func respondError(rw http.ResponseWriter, req *http.Request, message string, code int) {
	respond(rw, req, pages.Error, struct {
		Message string
		Code    int
	}{
		Message: message,
		Code:    code,
	}, code)
}

func waitForCallback(fc *connector.FakeConnector, ID manifold.ID) bool {
	timeout := time.After(callbackTimeout)
waitForCallback:
	select {
	case cb := <-fc.OnCallback:
		if cb.ID != ID {
			goto waitForCallback
		} else if cb.State == connector.PendingCallbackState {
			goto waitForCallback
		}
		return true
	case <-timeout:
		return false
	}
}
