package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	manifold "github.com/manifoldco/go-manifold"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// CreateUser create user
// swagger:model CreateUser
type CreateUser struct {

	// body
	// Required: true
	Body *CreateUserBody `json:"body"`
}

// Validate validates this create user
func (m *CreateUser) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBody(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateUser) validateBody(formats strfmt.Registry) error {

	if err := validate.Required("body", "body", m.Body); err != nil {
		return err
	}

	if m.Body != nil {

		if err := m.Body.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body")
			}
			return err
		}
	}

	return nil
}

// CreateUserBody create user body
// swagger:model CreateUserBody
type CreateUserBody struct {

	// email
	// Required: true
	Email manifold.Email `json:"email"`

	// name
	// Required: true
	Name UserDisplayName `json:"name"`

	// public key
	// Required: true
	PublicKey *LoginPublicKey `json:"public_key"`
}

// Validate validates this create user body
func (m *CreateUserBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEmail(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validatePublicKey(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateUserBody) validateEmail(formats strfmt.Registry) error {

	if err := m.Email.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("body" + "." + "email")
		}
		return err
	}

	return nil
}

func (m *CreateUserBody) validateName(formats strfmt.Registry) error {

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("body" + "." + "name")
		}
		return err
	}

	return nil
}

func (m *CreateUserBody) validatePublicKey(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"public_key", "body", m.PublicKey); err != nil {
		return err
	}

	if m.PublicKey != nil {

		if err := m.PublicKey.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "public_key")
			}
			return err
		}
	}

	return nil
}
