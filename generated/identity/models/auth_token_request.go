package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// AuthTokenRequest auth token request
// swagger:model AuthTokenRequest
type AuthTokenRequest struct {

	// login token sig
	// Required: true
	// Max Length: 86
	// Min Length: 86
	// Pattern: ^[a-zA-Z0-9_-]*$
	LoginTokenSig *string `json:"login_token_sig"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this auth token request
func (m *AuthTokenRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLoginTokenSig(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AuthTokenRequest) validateLoginTokenSig(formats strfmt.Registry) error {

	if err := validate.Required("login_token_sig", "body", m.LoginTokenSig); err != nil {
		return err
	}

	if err := validate.MinLength("login_token_sig", "body", string(*m.LoginTokenSig), 86); err != nil {
		return err
	}

	if err := validate.MaxLength("login_token_sig", "body", string(*m.LoginTokenSig), 86); err != nil {
		return err
	}

	if err := validate.Pattern("login_token_sig", "body", string(*m.LoginTokenSig), `^[a-zA-Z0-9_-]*$`); err != nil {
		return err
	}

	return nil
}

var authTokenRequestTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["auth"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		authTokenRequestTypeTypePropEnum = append(authTokenRequestTypeTypePropEnum, v)
	}
}

const (
	// AuthTokenRequestTypeAuth captures enum value "auth"
	AuthTokenRequestTypeAuth string = "auth"
)

// prop value enum
func (m *AuthTokenRequest) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, authTokenRequestTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *AuthTokenRequest) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}
