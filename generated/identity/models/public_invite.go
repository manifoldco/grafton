package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	manifold "github.com/manifoldco/go-manifold"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PublicInvite public invite
// swagger:model PublicInvite
type PublicInvite struct {

	// invited
	// Required: true
	Invited *PublicInviteInvited `json:"invited"`

	// invited by
	// Required: true
	InvitedBy *PublicInviteInvitedBy `json:"invited_by"`

	// role
	// Required: true
	Role RoleLabel `json:"role"`

	// team
	// Required: true
	Team *Team `json:"team"`
}

// Validate validates this public invite
func (m *PublicInvite) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInvited(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateInvitedBy(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateRole(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTeam(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicInvite) validateInvited(formats strfmt.Registry) error {

	if err := validate.Required("invited", "body", m.Invited); err != nil {
		return err
	}

	if m.Invited != nil {

		if err := m.Invited.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("invited")
			}
			return err
		}
	}

	return nil
}

func (m *PublicInvite) validateInvitedBy(formats strfmt.Registry) error {

	if err := validate.Required("invited_by", "body", m.InvitedBy); err != nil {
		return err
	}

	if m.InvitedBy != nil {

		if err := m.InvitedBy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("invited_by")
			}
			return err
		}
	}

	return nil
}

func (m *PublicInvite) validateRole(formats strfmt.Registry) error {

	if err := m.Role.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("role")
		}
		return err
	}

	return nil
}

func (m *PublicInvite) validateTeam(formats strfmt.Registry) error {

	if err := validate.Required("team", "body", m.Team); err != nil {
		return err
	}

	if m.Team != nil {

		if err := m.Team.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("team")
			}
			return err
		}
	}

	return nil
}

// PublicInviteInvited public invite invited
// swagger:model PublicInviteInvited
type PublicInviteInvited struct {

	// email
	Email manifold.Email `json:"email,omitempty"`

	// name
	Name UserDisplayName `json:"name,omitempty"`
}

// Validate validates this public invite invited
func (m *PublicInviteInvited) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicInviteInvited) validateName(formats strfmt.Registry) error {

	if swag.IsZero(m.Name) { // not required
		return nil
	}

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("invited" + "." + "name")
		}
		return err
	}

	return nil
}

// PublicInviteInvitedBy public invite invited by
// swagger:model PublicInviteInvitedBy
type PublicInviteInvitedBy struct {

	// email
	Email manifold.Email `json:"email,omitempty"`

	// name
	Name UserDisplayName `json:"name,omitempty"`
}

// Validate validates this public invite invited by
func (m *PublicInviteInvitedBy) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicInviteInvitedBy) validateName(formats strfmt.Registry) error {

	if swag.IsZero(m.Name) { // not required
		return nil
	}

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("invited_by" + "." + "name")
		}
		return err
	}

	return nil
}
