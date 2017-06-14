package acceptance

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/grafton/connector"
)

var sso = Feature("sso", "Single Sign-On Flow", func(ctx context.Context) {
	Default(func() {
		authCode, err := fakeConnector.CreateCode()
		if err != nil {
			FatalErr("could not create auth code", err)
		}

		url := api.CreateSsoURL(authCode.Code, resourceID)
		Infoln("Attempting to SSO into URL:", url)

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			FatalErr("got error building new request", err)
		}

		client := http.Client{
			// don't follow redirects.
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req = req.WithContext(ctx)
		resp, err := client.Do(req)

		logRequest(req)
		logResponse(resp)

		gm.Expect(err).To(notError())

		capturer, err := fakeConnector.GetCapturer("/v1/oauth/tokens")
		if err != nil {
			FatalErr("Could not find request capturer")
		}

		foundReqs := capturer.Get()
		reqs := []interface{}{}
		for _, v := range foundReqs {
			req, ok := v.(*connector.TokenRequest)
			if !ok {
				FatalErr("Could not cast request body to TokenRequest")
			}

			if req.GrantType == connector.AuthorizationCodeGrantType {
				reqs = append(reqs, req)
			}
		}

		gm.Expect(resp.StatusCode).To(gm.SatisfyAny(
			gm.BeNumerically("==", 200),
			gm.BeNumerically("==", 302),
			gm.BeNumerically("==", 303),
		), "Status code should be success (200) or redirect (302 or 303)")

		gm.Expect(len(reqs)).To(
			gm.Equal(1), "Zero or more than one token request should be received")

		tokReq := reqs[0].(*connector.TokenRequest)
		gm.Expect(tokReq).To(matchTokenRequest(&connector.TokenRequest{
			ContentType:  "application/x-www-form-urlencoded",
			GrantType:    "authorization_code",
			Code:         authCode.Code,
			AuthHeader:   tokReq.AuthHeader,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}), "Invalid token request")
	})

	ErrorCase("with wrong client id", func() {
		fakeConnector.Config.ClientID = "fake-client"
		defer func() {
			fakeConnector.Config.ClientID = clientID
		}()
		authCode, err := fakeConnector.CreateCode()
		if err != nil {
			FatalErr("could not create auth code", err)
		}

		url := api.CreateSsoURL(authCode.Code, resourceID)
		Infoln("Attempting to SSO into URL:", url)

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			FatalErr("got error building new request", err)
		}

		client := http.Client{
			// don't follow redirects.
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req = req.WithContext(ctx)
		resp, err := client.Do(req)

		logRequest(req)
		logResponse(resp)

		gm.Expect(err).To(notError())
		defer resp.Body.Close()

		gm.Expect(resp.StatusCode).To(gm.SatisfyAll(
			gm.BeNumerically("==", 401),
		), "Status code should be 401 unauthorized")
	})

	ErrorCase("with wrong client secret", func() {
		fakeConnector.Config.ClientSecret = "fake-secret"
		defer func() {
			fakeConnector.Config.ClientSecret = clientSecret
		}()

		authCode, err := fakeConnector.CreateCode()
		if err != nil {
			FatalErr("could not create auth code", err)
		}

		url := api.CreateSsoURL(authCode.Code, resourceID)
		Infoln("Attempting to SSO into URL:", url)

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			FatalErr("got error building new request", err)
		}

		client := http.Client{
			// don't follow redirects.
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req = req.WithContext(ctx)
		resp, err := client.Do(req)

		logRequest(req)
		logResponse(resp)

		gm.Expect(err).To(notError())
		defer resp.Body.Close()

		gm.Expect(resp.StatusCode).To(gm.SatisfyAll(
			gm.BeNumerically("==", 401),
		), "Status code should be 401 unauthorized")
	})

	ErrorCase("with expired token", func() {
		authCode, err := fakeConnector.CreateCode()
		authCode.ExpiresAt = time.Now().Add(-1 * time.Minute)
		if err != nil {
			FatalErr("could not create auth code", err)
		}

		url := api.CreateSsoURL(authCode.Code, resourceID)
		Infoln("Attempting to SSO into URL:", url)

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			FatalErr("got error building new request", err)
		}

		client := http.Client{
			// don't follow redirects.
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req = req.WithContext(ctx)
		resp, err := client.Do(req)

		logRequest(req)
		logResponse(resp)

		gm.Expect(err).To(notError())
		defer resp.Body.Close()

		gm.Expect(resp.StatusCode).To(gm.SatisfyAll(
			gm.BeNumerically("==", 401),
		), "Status code should be 401 unauthorized")
	})

	ErrorCase("with non-existing code", func() {
		if _, err := fakeConnector.CreateCode(); err != nil {
			FatalErr("could not create auth code", err)
		}

		wrongCode := &connector.AuthorizationCode{
			Code:      "non-existing",
			ExpiresAt: time.Now().Add(3600 * time.Second),
		}

		url := uapi.CreateSsoURL(wrongCode.Code, resourceID)
		Infoln("Attempting to SSO into URL:", url)

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			FatalErr("got error building new request", err)
		}

		client := http.Client{
			// don't follow redirects.
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req = req.WithContext(ctx)
		resp, err := client.Do(req)

		logRequest(req)
		logResponse(resp)

		gm.Expect(err).To(notError())
		defer resp.Body.Close()

		gm.Expect(resp.StatusCode).To(gm.SatisfyAll(
			gm.BeNumerically("==", 401),
		), "Status code should be 401 unauthorized")
	})

	ErrorCase("with connector response error", func() {
		fakeConnector.Server.Handler = connector.ErrorHandler(fakeConnector)
		defer func() {
			fakeConnector.Server.Handler = connector.ValidHandler(fakeConnector)
		}()

		authCode, err := fakeConnector.CreateCode()
		if err != nil {
			FatalErr("could not create auth code", err)
		}

		url := api.CreateSsoURL(authCode.Code, resourceID)
		Infoln("Attempting to SSO into URL:", url)

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			FatalErr("got error building new request", err)
		}

		client := http.Client{
			// don't follow redirects.
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req = req.WithContext(ctx)
		resp, err := client.Do(req)

		logRequest(req)
		logResponse(resp)

		gm.Expect(err).To(notError())
		defer resp.Body.Close()

		gm.Expect(resp.StatusCode).To(gm.SatisfyAll(
			gm.BeNumerically("==", 401),
		), "Status code should be 401 unauthorized")
	})
})

var _ = sso.RunsInside("provision")
var _ = sso.RequiredFlags("client-id", "client-secret", "connector-port")

func matchTokenRequest(expected *connector.TokenRequest) *expectedTokenRequestMatcher {
	return &expectedTokenRequestMatcher{expected: expected}
}

type expectedTokenRequestMatcher struct {
	expected *connector.TokenRequest
}

func (e expectedTokenRequestMatcher) Match(actual interface{}) (bool, error) {
	return reflect.DeepEqual(e.expected, actual), nil

}

func (e expectedTokenRequestMatcher) FailureMessage(actual interface{}) string {
	other := actual.(*connector.TokenRequest)

	if other == nil {
		return "Did not receive an OAuth2 token request"
	}

	msg := ""
	if e.expected.ContentType != other.ContentType {
		msg += fmt.Sprintf("Expected content-type of '%s', but got '%s'.\n",
			e.expected.ContentType, other.ContentType)
	}

	source := ""
	if other.AuthHeader {
		source = " (read from Authorization header)"
	}

	if e.expected.ClientID != other.ClientID {
		msg += fmt.Sprintf("Expected client_id of '%s', but got '%s'%s.\n",
			e.expected.ClientID, other.ClientID, source)
	}
	if e.expected.ClientSecret != other.ClientSecret {
		msg += fmt.Sprintf("Expected client_secret of '%s', but got '%s'%s.\n",
			e.expected.ClientSecret, other.ClientSecret, source)
	}

	if e.expected.Code != other.Code {
		msg += fmt.Sprintf("Expected code of '%s', but got '%s'.\n",
			e.expected.Code, other.Code)
	}
	if e.expected.GrantType != other.GrantType {
		msg += fmt.Sprintf("Expected grant_type of '%s', but got '%s'.\n",
			e.expected.GrantType, other.GrantType)
	}

	return msg
}

func (e expectedTokenRequestMatcher) NegatedFailureMessage(actual interface{}) string {
	// We don't use the negated form.
	return "Token should not have matched expected values"
}

func logRequest(req *http.Request) {
	rq, _ := httputil.DumpRequest(req, true)
	Infoln(string(rq))
}

func logResponse(rsp *http.Response) {
	resp, _ := httputil.DumpResponse(rsp, true)
	Infoln(string(resp))
}
