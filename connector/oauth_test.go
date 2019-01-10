package connector

import (
	"encoding/base64"
	gm "github.com/onsi/gomega"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
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
			req.Header.Add("Content-Type", "application/json")
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

	t.Run("create access token [client_credentials] responds properly to a body and a Authorization header",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(`{
				"grant_type": "client_credentials"
			}`))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(201))
		})

	t.Run("create access token [authorization_code] responds properly to json data",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			authCode, err := c.CreateCode()
			gm.Expect(err).ToNot(gm.HaveOccurred())

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(`{
				"grant_type": "authorization_code",
				"client_id": "`+clientID+`",
				"client_secret": "`+clientSecret+`",
				"code": "`+authCode.Code+`"
			}`))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(201))
		})

	t.Run("create access token [authorization_code] responds properly to url encoded data",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			authCode, err := c.CreateCode()
			gm.Expect(err).ToNot(gm.HaveOccurred())

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(``))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.PostForm = make(url.Values, 3)
			req.PostForm.Add("grant_type", "authorization_code")
			req.PostForm.Add("client_id", clientID)
			req.PostForm.Add("client_secret", clientSecret)
			req.PostForm.Add("code", authCode.Code)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(201))
		})

	t.Run("create access token [authorization_code] responds properly to json data and an Authorization header",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			authCode, err := c.CreateCode()
			gm.Expect(err).ToNot(gm.HaveOccurred())

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(`{
				"grant_type": "authorization_code",
				"code": "`+authCode.Code+`"
			}`))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(201))
		})

	t.Run("create access token will return an error for an empty or invalid request body",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(""))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(400))

			req = httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader("{"))
			req.Header.Add("Content-Type", "application/json")
			rec = httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(400))
		})

	t.Run("create access token will return an error for an invalid grant_type",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(`{
				"grant_type": "random",
				"client_id": "`+clientID+`",
				"client_secret": "`+clientSecret+`"
			}`))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(400))
		})

	t.Run("create access token [authorization_code] will return an error if the code is missing",
		func(t *testing.T) {
			gm.RegisterTestingT(t)

			req := httptest.NewRequest("POST", "/oauth/tokens", strings.NewReader(`{
				"grant_type": "authorization_code",
				"client_id": "`+clientID+`",
				"client_secret": "`+clientSecret+`"
			}`))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			gm.Expect(rec.Code).To(gm.Equal(400))
		})
}
