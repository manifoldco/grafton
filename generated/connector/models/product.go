package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/manifoldco/go-manifold"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// Product product
// swagger:model product
type Product struct {

	// target
	// Required: true
	Target *ProductAO1Target `json:"target"`
}

func (m *Product) Type() string {
	return "product"
}
func (m *Product) SetType(val string) {

}

// UnmarshalJSON unmarshals this polymorphic type from a JSON structure
func (m *Product) UnmarshalJSON(raw []byte) error {
	var data struct {
		Type string `json:"type"`

		// target
		// Required: true
		Target *ProductAO1Target `json:"target"`
	}

	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	m.Target = data.Target

	return nil
}

// MarshalJSON marshals this polymorphic type to a JSON structure
func (m Product) MarshalJSON() ([]byte, error) {
	var data struct {
		Type string `json:"type"`

		// target
		// Required: true
		Target *ProductAO1Target `json:"target"`
	}

	data.Target = m.Target
	data.Type = "product"
	return json.Marshal(data)
}

// Validate validates this product
func (m *Product) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTarget(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Product) validateTarget(formats strfmt.Registry) error {

	if err := validate.Required("target", "body", m.Target); err != nil {
		return err
	}

	if m.Target != nil {

		if err := m.Target.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("target")
			}
			return err
		}
	}

	return nil
}

// ProductAO1Target product a o1 target
// swagger:model ProductAO1Target
type ProductAO1Target struct {

	// label
	// Required: true
	Label manifold.Label `json:"label"`

	// name
	// Required: true
	Name manifold.Name `json:"name"`
}

// Validate validates this product a o1 target
func (m *ProductAO1Target) Validate(formats strfmt.Registry) error {
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

func (m *ProductAO1Target) validateLabel(formats strfmt.Registry) error {

	if err := m.Label.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("target" + "." + "label")
		}
		return err
	}

	return nil
}

func (m *ProductAO1Target) validateName(formats strfmt.Registry) error {

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("target" + "." + "name")
		}
		return err
	}

	return nil
}
