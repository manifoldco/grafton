package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
)

// ResourceSettings **ALPHA** Warning: this API may change in the future or completely removed. Do NOT used for
// product system.
// Object describing additional settings for the resource specified by the user during provision.
//
// swagger:model ResourceSettings
type ResourceSettings map[string]string

// Validate validates this resource settings
func (m ResourceSettings) Validate(formats strfmt.Registry) error {
	return nil
}
