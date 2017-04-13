package connector

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/runtime"

	"github.com/manifoldco/go-connector"
	cerrors "github.com/manifoldco/go-connector/errors"
	"github.com/manifoldco/go-jwt"
	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"

	"github.com/manifoldco/grafton"
)

// Errors for authorizing requests
var (
	errMissingAuthHeader  = grafton.NewError(errors.BadRequestError, "Missing Authorization Header")
	errInvalidAuthHeader  = grafton.NewError(errors.BadRequestError, "Invalid Authorization Header")
	errInvalidAccessToken = grafton.NewError(errors.BadRequestError, "Invalid access token")
	errUnauthorized       = grafton.NewError(errors.UnauthorizedError, "Unauthorized")
)

// Errors for oauth flow
var (
	errUnsupportedGrantType = connector.NewOAuthError(cerrors.UnsupportedGrantErrorType, "Unsupported grant type")

	errInvalidOAuthContentType = connector.NewOAuthError(cerrors.InvalidRequestErrorType, "Invalid content type")
	errInvalidClientCreds      = connector.NewOAuthError(cerrors.InvalidClientErrorType, "Invalid client credentials")

	errMissingCode = connector.NewOAuthError(cerrors.InvalidGrantErrorType, "No code provided")
	errExpiredCode = connector.NewOAuthError(cerrors.InvalidGrantErrorType, "Authorization code has expired")
)

type claims struct {
	ClientID string
	TokenID  manifold.ID
}

var jsonProducer = runtime.JSONProducer()

func respondWithError(rw http.ResponseWriter, err error) {
	e := manifold.ToError(err)
	e.WriteResponse(rw, jsonProducer)
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

func authorizeRequest(c *FakeConnector, req *http.Request) (*AccessToken, error) {
	h := req.Header.Get("Authorization")
	if h == "" {
		return nil, errMissingAuthHeader
	}

	parts := strings.Split(h, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errInvalidAuthHeader
	}

	_, err := jwt.Read(c.Config.SigningKey, parts[1])
	if err != nil {
		return nil, errInvalidAccessToken
	}

	token := c.getToken(parts[1])
	if token == nil {
		return nil, errUnauthorized
	}

	return token, nil
}

func getSelfHandler(c *FakeConnector, capturer *RequestCapturer) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		token, err := authorizeRequest(c, req)
		if err != nil {
			respondWithError(rw, err)
			return
		}

		var body interface{}
		switch token.GrantType {
		case AuthorizationCodeGrantType:
			body = UserProfile{
				Type: "user",
				Target: &UserTarget{
					Name:  "joe user",
					Email: "joe@user.com",
				},
			}
		case ClientCredentialsGrantType:
			body = ProductProfile{
				Type: "product",
				Target: &ProductTarget{
					Name:  "A Great Product",
					Label: c.Config.Product,
				},
			}
		default:
			respondWithError(rw, errISE)
			return
		}

		respondWithJSON(rw, body, 200)
	}
}

func createAccessTokenHandler(c *FakeConnector, capturer *RequestCapturer) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		id, secret, ok := req.BasicAuth()
		if !ok {
			id = req.FormValue("client_id")
			secret = req.FormValue("client_secret")
		}

		var gt GrantType
		switch req.FormValue("grant_type") {
		case "authorization_code":
			gt = AuthorizationCodeGrantType
		case "client_credentials":
			gt = ClientCredentialsGrantType
		default:
			// We'll allow this for now, to capture it in the request recording,
			// but respond with an error in the switch below.
			gt = GrantType(req.FormValue("grant_type"))
		}

		tokReq := &TokenRequest{
			ContentType:  req.Header.Get("Content-Type"),
			Code:         req.FormValue("code"),
			AuthHeader:   ok,
			ClientID:     id,
			ClientSecret: secret,
			GrantType:    gt,
		}

		capturer.capture(tokReq)

		var e *connector.OAuthError
		switch tokReq.GrantType {
		case AuthorizationCodeGrantType:
			e = validateAuthCodeGrant(c, tokReq)
		case ClientCredentialsGrantType:
			e = validateClientCredentialGrant(c, tokReq)
		default:
			e = errUnsupportedGrantType
		}

		if e != nil {
			e.WriteResponse(rw, jsonProducer)
			return
		}

		tokenID, err := manifold.NewID(idtype.OAuthAccessToken)
		if err != nil {
			e = connector.ToOAuthError(err).(*connector.OAuthError)
			e.WriteResponse(rw, jsonProducer)
			return
		}

		jwtString, _, err := jwt.New(c.Config.SigningKey, &claims{
			ClientID: c.Config.ClientID,
			TokenID:  tokenID,
		}, nil)
		if err != nil {
			e = connector.ToOAuthError(err).(*connector.OAuthError)
			e.WriteResponse(rw, jsonProducer)
			return
		}

		t := &AccessToken{
			AccessToken: jwtString,
			ExpiresIn:   3600,
			TokenType:   "bearer",
			GrantType:   tokReq.GrantType,
			ID:          tokenID,
		}

		c.addToken(t)
		respondWithJSON(rw, t, 200)
	}
}

func validateAuthCodeGrant(c *FakeConnector, t *TokenRequest) *connector.OAuthError {
	err := validateClientCredentialGrant(c, t)
	if err != nil {
		return err
	}

	switch {
	case c.getCode(t.Code) == nil:
		err = errMissingCode
	case c.getCode(t.Code).ExpiresAt.Unix()-time.Now().UTC().Unix() < 1:
		err = errExpiredCode
	}

	return err
}

func validateClientCredentialGrant(c *FakeConnector, t *TokenRequest) *connector.OAuthError {
	var err *connector.OAuthError
	switch {
	case t.ContentType != "application/x-www-form-urlencoded":
		err = errInvalidOAuthContentType
	case t.ClientID != c.Config.ClientID:
		fallthrough
	case t.ClientSecret != c.Config.ClientSecret:
		err = errInvalidClientCreds
	}

	return err
}
