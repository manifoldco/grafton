package connector

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestCreateAccessTokenHandler(t *testing.T) {
	gm.RegisterTestingT(t)

	c := getConnectorInstance()
	r := &RequestCapturer{
		Route:    "/oauth/tokens",
		requests: make([]interface{}, 0),
	}
	handler := createAccessTokenHandler(c, r)

	t.Run("create access token [client_credentials] responds properly to json data",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(`{
				"grant_type": "client_credentials",
				"client_id": "`+clientID+`",
				"client_secret": "`+clientSecret+`"
			}`))
			req.Header.Add("Content-Type", "text/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(201))
		})

	t.Run("create access token [client_credentials] responds properly to url encoded data",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(``))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.PostForm = make(url.Values, 3)
			req.PostForm.Add("grant_type", "client_credentials")
			req.PostForm.Add("client_id", clientID)
			req.PostForm.Add("client_secret", clientSecret)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(201))
		})
}
