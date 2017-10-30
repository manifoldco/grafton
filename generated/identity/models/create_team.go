package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	manifold "github.com/manifoldco/go-manifold"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// CreateTeam create team
// swagger:model CreateTeam
type CreateTeam struct {

	// body
	// Required: true
	Body *CreateTeamBody `json:"body"`
}

// Validate validates this create team
func (m *CreateTeam) Validate(formats strfmt.Registry) error {
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

func (m *CreateTeam) validateBody(formats strfmt.Registry) error {

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

// CreateTeamBody create team body
// swagger:model CreateTeamBody
type CreateTeamBody struct {

	// label
	// Required: true
	Label manifold.Label `json:"label"`

	// name
	// Required: true
	Name manifold.Name `json:"name"`
}

// Validate validates this create team body
func (m *CreateTeamBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLabel(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateTeamBody) validateLabel(formats strfmt.Registry) error {

	if err := m.Label.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("body" + "." + "label")
		}
		return err
	}

	return nil
}

func (m *CreateTeamBody) validateName(formats strfmt.Registry) error {

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("body" + "." + "name")
		}
		return err
	}

	return nil
}
