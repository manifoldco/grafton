package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/packr"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/grafton/connector"
)

const callbackTimeout = time.Minute * 5

var templates packr.Box

func init() {
	templates = packr.NewBox("../templates")
}

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

func respondWithHTML(rw http.ResponseWriter, name string, v interface{}, code int) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(code)

	if v == nil {
		return
	}

	tpl := template.New("layout.html")

	// Parse in Route Page
	pageHTML, err := templates.FindString(name + ".html")
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}
	tpl, err = tpl.Parse(pageHTML)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// Parse in Layout
	layoutHTML, err := templates.FindString("layout.html")
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}
	tpl, err = tpl.Parse(layoutHTML)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	err = tpl.Execute(rw, v)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}
}

func respond(rw http.ResponseWriter, req *http.Request, templatePath string, v interface{}, code int) {
	// If JSON Request, respond as such
	if isJSONRequest(req) {
		respondWithJSON(rw, v, code)
		return
	}
	// Else HTML
	respondWithHTML(rw, templatePath, v, code)
}

func respondError(rw http.ResponseWriter, req *http.Request, message string, code int) {
	respond(rw, req, "error", struct {
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
