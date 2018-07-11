package primitives

import (
	"time"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/idtype"
)

// OAuthAuthorizationCodeBody represents the contents of an Authorization Code
type OAuthAuthorizationCodeBody struct {
	UserID            manifold.ID  `json:"user_id"`
	TeamID            *manifold.ID `json:"team_id,omitempty"`
	InboundResourceID manifold.ID  `json:"resource_id"`
	CreatedAt         time.Time    `json:"created_at"`
	ExpiresAt         time.Time    `json:"expires_at"`
	RedirectURI       string       `json:"redirect_uri,omitempty"`
	Code              string       `json:"code"`
}

// OAuthAuthorizationCode represents a short-term code used by a provider to create
// an Access Token scoped to a User and Resource
type OAuthAuthorizationCode struct {
	ID            manifold.ID                `json:"id"`
	StructType    string                     `json:"type"`
	StructVersion int                        `json:"version"`
	Body          OAuthAuthorizationCodeBody `json:"body"`
}

// GetID returns the current Identifier for this object
func (a *OAuthAuthorizationCode) GetID() manifold.ID {
	return a.ID
}

// Version returns the current struct version for this object
func (a *OAuthAuthorizationCode) Version() int {
	return a.StructVersion
}

// Type returns the OAuthAuthorizationCode idtype
func (a *OAuthAuthorizationCode) Type() idtype.Type {
	return idtype.OAuthAuthorizationCode
}

// Mutable represents the fact that this is a Mutable struct
func (a *OAuthAuthorizationCode) Mutable() {}
