package connector

import (
	"sync"
	"time"

	"github.com/manifoldco/go-manifold"
)

// GrantType represents a type of access token grant
type GrantType string

// AuthorizationCodeGrantType represents an OAuth Authorization Code
var (
	AuthorizationCodeGrantType GrantType = "authorization_code"
	ClientCredentialsGrantType GrantType = "client_credentials"
)

// TokenRequest represents all of the important values from a request by a
// provider to create an OAuth Token
type TokenRequest struct {
	ContentType  string
	Code         string
	AuthHeader   bool
	ClientID     string
	ClientSecret string
	GrantType    GrantType
	ExpiresAt    time.Time
}

// AuthorizationCode represents a code granted by the fake connector for
// kicking off the oauth flow
type AuthorizationCode struct {
	Code      string
	ExpiresAt time.Time
}

// AccessToken represents an access token granted by the fake connector for
// kicking off the oauth flow
type AccessToken struct {
	ID          manifold.ID `json:"-"`
	AccessToken string      `json:"access_token"`
	ExpiresIn   int         `json:"expires_in"`
	TokenType   string      `json:"token_type"`
	GrantType   GrantType   `json:"-"`
}

// Resource represents a resource provisioned through Grafton
type Resource struct {
	ID        manifold.ID `json:"id"`
	Plan      string      `json:"plan"`
	Product   string      `json:"product"`
	Region    string      `json:"region"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// UserProfile represents the data returned on GET /v1/self when the target
// type is a user
type UserProfile struct {
	Type   string      `json:"type"`
	Target *UserTarget `json:"target"`
}

// UserTarget represents the contents of a response on GET /v1/self when the
// target is a user
type UserTarget struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ProductTarget represents the data returned on GET /v1/vself when the target
// type is product
type ProductTarget struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

// ProductProfile represents the data returned on GET /v1/self when the target
// type is a product
type ProductProfile struct {
	Type   string         `json:"type"`
	Target *ProductTarget `json:"target"`
}

// CallbackType represents a type of callback
type CallbackType string

// ResourceProvisionCallback represents a callback related to a resource
// provision
var (
	ResourceProvisionCallback     CallbackType = "resource:provision"
	CredentialProvisionCallback   CallbackType = "credential:provision"
	ResourceDeprovisionCallback   CallbackType = "resource:deprovision"
	CredentialDeprovisionCallback CallbackType = "credential:deprovision"
	ResourceResizeCallback        CallbackType = "resource:resize"
)

// CallbackState represents a state reported by the provider regarding an
// operation
type CallbackState string

// PendingCallbackState represents a callback which is pending (not resolved)
var (
	PendingCallbackState CallbackState = "pending"
	DoneCallbackState    CallbackState = "done"
	ErrorCallbackState   CallbackState = "error"
)

// Callback represents a callback that is either pending or has been received
// from a provider
type Callback struct {
	ID          manifold.ID       `json:"id"`
	Mutex       *sync.Mutex       `json:"-"`
	Type        CallbackType      `json:"type"`
	State       CallbackState     `json:"state"`
	Message     string            `json:"message"`
	Credentials map[string]string `json:"-"`
}

// CallbackRequest represents a received callback from a provider
type CallbackRequest struct {
	State       CallbackState     `json:"state"`
	Message     string            `json:"message"`
	Credentials map[string]string `json:"credentials"`
}
