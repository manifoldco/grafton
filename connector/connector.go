package connector

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-zoo/bone"

	"github.com/manifoldco/go-base32"
	"github.com/manifoldco/go-connector"
	cerrors "github.com/manifoldco/go-connector/errors"
	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/idtype"
	"github.com/manifoldco/grafton/db"
)

// ErrCallbackNotFound represents an error which occurs if a callback does not
// exist
var ErrCallbackNotFound = errors.New("Callback Not Found")

// ErrCallbackAlreadyResolved represents an error which occurs if the callback
// has already been resolved
var ErrCallbackAlreadyResolved = errors.New("Callback Already Resolved")

// ErrResourceNotFound represents an error which occurrs if the resource does
// not exist
var ErrResourceNotFound = errors.New("Resource Not Found")

// RequestCapturer represents functionality for capturing and storing requests
// for a specific route
type RequestCapturer struct {
	Route    string
	requests []interface{}
}

// capture holds onto a captured requests
func (r *RequestCapturer) capture(v interface{}) {
	r.requests = append(r.requests, v)
}

// Get returns the requests captured by this Capturer
func (r *RequestCapturer) Get() []interface{} {
	return r.requests
}

// FakeConnectorConfig represents the values used to Configure a FakeConnector
type FakeConnectorConfig struct {
	Product      string
	Port         uint
	ClientID     string
	ClientSecret string
	SigningKey   string
}

// FakeConnector represents a fake connector api server run by Grafton for use
// by providers to integrate with Manifold.
type FakeConnector struct {
	Config     *FakeConnectorConfig
	DB         *db.DB
	OnCallback chan *Callback
	capturers  map[string]*RequestCapturer
	codes      []*AuthorizationCode
	tokens     []*AccessToken
	callbacks  []*Callback
	Server     *http.Server
}

// StartSync starts the server or returns an error if it couldn't be started
func (c *FakeConnector) StartSync() error {
	h := ValidHandler(c)
	c.Server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", c.Config.Port),
		Handler: h,
	}

	return c.Server.ListenAndServe()
}

// Start the server or return an error if it couldn't be started
func (c *FakeConnector) Start() {
	h := ValidHandler(c)
	c.Server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", c.Config.Port),
		Handler: h,
	}

	go c.Server.ListenAndServe()
}

// Stop the server or return an error if it couldn't be stopped
func (c *FakeConnector) Stop() error {
	if c.Server == nil {
		return errors.New("Cannot not stop a server that has not started")
	}

	close(c.OnCallback)
	return c.Server.Close()
}

// GetCapturer returns a RequestCapturer for the given route, if no capturer
// exists, an error is returned.
func (c *FakeConnector) GetCapturer(route string) (*RequestCapturer, error) {
	capturer, ok := c.capturers[route]
	if !ok {
		return nil, errors.New("No Capturer found")
	}

	return capturer, nil
}

// AddResource stores a resource inside the connector
func (c *FakeConnector) AddResource(r *db.Resource) {
	c.DB.PutResource(*r)
}

// RemoveResource deletes a resource stored inside the connector
func (c *FakeConnector) RemoveResource(ID manifold.ID) error {
	if c.DB.DeleteResource(ID) {
		return nil
	}
	return ErrResourceNotFound
}

// GetResource returns a Resource for the ID or nil
func (c *FakeConnector) GetResource(id manifold.ID) *db.Resource {
	return c.DB.GetResource(id)
}

// AddCallback stores the callback inside the FakeConnector
func (c *FakeConnector) AddCallback(t CallbackType) (*Callback, error) {
	ID, err := manifold.NewID(idtype.Callback)
	if err != nil {
		return nil, err
	}

	cb := &Callback{
		ID:          ID,
		Mutex:       &sync.Mutex{},
		State:       PendingCallbackState,
		Type:        t,
		Message:     "",
		Credentials: make(map[string]string),
	}

	c.callbacks = append(c.callbacks, cb)

	return cb, nil
}

// GetCallback returns a callback for the given id if it exists
func (c *FakeConnector) GetCallback(ID manifold.ID) *Callback {
	for _, v := range c.callbacks {
		if v.ID == ID {
			return v
		}
	}

	return nil
}

// TriggerCallback updates the callback if it's still pending, and notifies any
// listeners
func (c *FakeConnector) TriggerCallback(ID manifold.ID, state CallbackState, msg string, creds map[string]string) error {
	cb := c.GetCallback(ID)
	if cb == nil {
		return ErrCallbackNotFound
	}

	cb.Mutex.Lock()
	defer cb.Mutex.Unlock()

	if cb.State != PendingCallbackState {
		if callbackEqual(cb, state, msg, creds) {
			return nil
		}

		return ErrCallbackAlreadyResolved
	}

	cb.State = state
	cb.Message = msg

	for k, v := range creds {
		cb.Credentials[k] = v
	}

	c.OnCallback <- cb
	return nil
}

func callbackEqual(cb *Callback, s CallbackState, msg string, creds map[string]string) bool {
	if cb.Message != msg || cb.State != s || len(creds) != len(cb.Credentials) {
		return false
	}

	for k := range creds {
		if creds[k] != cb.Credentials[k] {
			return false
		}
	}

	return true
}

func (c *FakeConnector) capturer(route string) *RequestCapturer {
	r := &RequestCapturer{
		Route:    route,
		requests: make([]interface{}, 0),
	}

	c.capturers[route] = r

	return r
}

func (c *FakeConnector) getToken(token string) *AccessToken {
	for _, v := range c.tokens {
		if v.AccessToken == token {
			return v
		}
	}

	return nil
}

func (c *FakeConnector) addToken(t *AccessToken) {
	c.tokens = append(c.tokens, t)
}

// CreateCode returns an AuthorizationCode method
func (c *FakeConnector) CreateCode() (*AuthorizationCode, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	code := base32.EncodeToString(b)
	authCode := &AuthorizationCode{
		Code:      code,
		ExpiresAt: time.Now().Add(3600 * time.Second),
	}

	c.codes = append(c.codes, authCode)
	return authCode, nil
}

func (c *FakeConnector) getCode(code string) *AuthorizationCode {
	for _, v := range c.codes {
		if v.Code == code {
			return v
		}
	}

	return nil
}

// New creates and configures a FakeConnector
func New(port uint, clientID string, clientSecret string, product string) (*FakeConnector, error) {
	c := &FakeConnector{
		Config: &FakeConnectorConfig{
			Product:      product,
			Port:         port,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			SigningKey:   "hello",
		},
		DB:         db.New(),
		OnCallback: make(chan *Callback, 100),
		capturers:  make(map[string]*RequestCapturer),
	}

	return c, nil
}

// ValidHandler returns a set of endpoints that are valid and in use by the
// production connector.
func ValidHandler(c *FakeConnector) *bone.Mux {
	mux := bone.New()
	mux.PostFunc("/v1/oauth/tokens", createAccessTokenHandler(c, c.capturer("/v1/oauth/tokens")))
	mux.GetFunc("/v1/self", getSelfHandler(c, c.capturer("/v1/self")))
	mux.PutFunc("/v1/callbacks/:id", processCallbackHandler(c, c.capturer("/v1/callbacks/{id}")))
	mux.GetFunc("/v1/resources/:id", getResourceHandler(c, c.capturer("/v1/resources/{id}")))
	mux.GetFunc("/v1/resources/:id/users", getResourceUsersHandler(c,
		c.capturer("/v1/resources/{id}/users")))
	mux.GetFunc("/v1/resources/:id/credentials", getResourceCredentialsHandler(c,
		c.capturer("/v1/resources/{id}/credentials")))
	mux.GetFunc("/v1/resources/:id/measures", getResourceMeasuresHandler(c,
		c.capturer("/v1/resources/{id}/measures")))
	mux.PutFunc("/v1/resources/:id/measures", putResourceMeasuresHandler(c,
		c.capturer("/v1/resources/{id}/measures")))
	return mux
}

// ErrorHandler is a set of endpoints that just returns errors. This is to test
// how the providers handle errors from our end.
func ErrorHandler(c *FakeConnector) *bone.Mux {
	mux := bone.New()
	mux.PostFunc("/v1/oauth/tokens", createErroredHandler(c, c.capturer("/v1/oauth/tokens")))
	return mux
}

func createErroredHandler(c *FakeConnector, capturer *RequestCapturer) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		e := connector.NewOAuthError(cerrors.ServerErrorErrorType, "Internal server error")
		e.WriteResponse(rw, jsonProducer)
	}
}
