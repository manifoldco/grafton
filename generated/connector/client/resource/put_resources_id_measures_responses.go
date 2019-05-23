package resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// PutResourcesIDMeasuresReader is a Reader for the PutResourcesIDMeasures structure.
type PutResourcesIDMeasuresReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutResourcesIDMeasuresReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 204:
		result := NewPutResourcesIDMeasuresNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewPutResourcesIDMeasuresBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewPutResourcesIDMeasuresUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewPutResourcesIDMeasuresNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewPutResourcesIDMeasuresInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPutResourcesIDMeasuresNoContent creates a PutResourcesIDMeasuresNoContent with default headers values
func NewPutResourcesIDMeasuresNoContent() *PutResourcesIDMeasuresNoContent {
	return &PutResourcesIDMeasuresNoContent{}
}

/*PutResourcesIDMeasuresNoContent handles this case with default header values.

Empty response
*/
type PutResourcesIDMeasuresNoContent struct {
}

func (o *PutResourcesIDMeasuresNoContent) Error() string {
	return fmt.Sprintf("[PUT /resources/{id}/measures][%d] putResourcesIdMeasuresNoContent ", 204)
}

func (o *PutResourcesIDMeasuresNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPutResourcesIDMeasuresBadRequest creates a PutResourcesIDMeasuresBadRequest with default headers values
func NewPutResourcesIDMeasuresBadRequest() *PutResourcesIDMeasuresBadRequest {
	return &PutResourcesIDMeasuresBadRequest{}
}

/*PutResourcesIDMeasuresBadRequest handles this case with default header values.

Request denied due to invalid request body, path, or headers.
*/
type PutResourcesIDMeasuresBadRequest struct {
	Payload PutResourcesIDMeasuresBadRequestBody
}

func (o *PutResourcesIDMeasuresBadRequest) Error() string {
	return fmt.Sprintf("[PUT /resources/{id}/measures][%d] putResourcesIdMeasuresBadRequest  %+v", 400, o.Payload)
}

func (o *PutResourcesIDMeasuresBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutResourcesIDMeasuresUnauthorized creates a PutResourcesIDMeasuresUnauthorized with default headers values
func NewPutResourcesIDMeasuresUnauthorized() *PutResourcesIDMeasuresUnauthorized {
	return &PutResourcesIDMeasuresUnauthorized{}
}

/*PutResourcesIDMeasuresUnauthorized handles this case with default header values.

Request denied as the provided credentials are no longer valid.
*/
type PutResourcesIDMeasuresUnauthorized struct {
	Payload PutResourcesIDMeasuresUnauthorizedBody
}

func (o *PutResourcesIDMeasuresUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /resources/{id}/measures][%d] putResourcesIdMeasuresUnauthorized  %+v", 401, o.Payload)
}

func (o *PutResourcesIDMeasuresUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutResourcesIDMeasuresNotFound creates a PutResourcesIDMeasuresNotFound with default headers values
func NewPutResourcesIDMeasuresNotFound() *PutResourcesIDMeasuresNotFound {
	return &PutResourcesIDMeasuresNotFound{}
}

/*PutResourcesIDMeasuresNotFound handles this case with default header values.

Request denied as the requested resource does not exist.
*/
type PutResourcesIDMeasuresNotFound struct {
	Payload PutResourcesIDMeasuresNotFoundBody
}

func (o *PutResourcesIDMeasuresNotFound) Error() string {
	return fmt.Sprintf("[PUT /resources/{id}/measures][%d] putResourcesIdMeasuresNotFound  %+v", 404, o.Payload)
}

func (o *PutResourcesIDMeasuresNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutResourcesIDMeasuresInternalServerError creates a PutResourcesIDMeasuresInternalServerError with default headers values
func NewPutResourcesIDMeasuresInternalServerError() *PutResourcesIDMeasuresInternalServerError {
	return &PutResourcesIDMeasuresInternalServerError{}
}

/*PutResourcesIDMeasuresInternalServerError handles this case with default header values.

Request failed due to an internal server error.
*/
type PutResourcesIDMeasuresInternalServerError struct {
	Payload PutResourcesIDMeasuresInternalServerErrorBody
}

func (o *PutResourcesIDMeasuresInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /resources/{id}/measures][%d] putResourcesIdMeasuresInternalServerError  %+v", 500, o.Payload)
}

func (o *PutResourcesIDMeasuresInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*PutResourcesIDMeasuresBadRequestBody put resources ID measures bad request body
swagger:model PutResourcesIDMeasuresBadRequestBody
*/
type PutResourcesIDMeasuresBadRequestBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this put resources ID measures bad request body
func (o *PutResourcesIDMeasuresBadRequestBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := o.validateType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutResourcesIDMeasuresBadRequestBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresBadRequest"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var putResourcesIdMeasuresBadRequestBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		putResourcesIdMeasuresBadRequestBodyTypeTypePropEnum = append(putResourcesIdMeasuresBadRequestBodyTypeTypePropEnum, v)
	}
}

const (
	// PutResourcesIDMeasuresBadRequestBodyTypeBadRequest captures enum value "bad_request"
	PutResourcesIDMeasuresBadRequestBodyTypeBadRequest string = "bad_request"
	// PutResourcesIDMeasuresBadRequestBodyTypeUnauthorized captures enum value "unauthorized"
	PutResourcesIDMeasuresBadRequestBodyTypeUnauthorized string = "unauthorized"
	// PutResourcesIDMeasuresBadRequestBodyTypeNotFound captures enum value "not_found"
	PutResourcesIDMeasuresBadRequestBodyTypeNotFound string = "not_found"
	// PutResourcesIDMeasuresBadRequestBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	PutResourcesIDMeasuresBadRequestBodyTypeMethodNotAllowed string = "method_not_allowed"
	// PutResourcesIDMeasuresBadRequestBodyTypeInternal captures enum value "internal"
	PutResourcesIDMeasuresBadRequestBodyTypeInternal string = "internal"
	// PutResourcesIDMeasuresBadRequestBodyTypeInvalidGrant captures enum value "invalid_grant"
	PutResourcesIDMeasuresBadRequestBodyTypeInvalidGrant string = "invalid_grant"
	// PutResourcesIDMeasuresBadRequestBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	PutResourcesIDMeasuresBadRequestBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *PutResourcesIDMeasuresBadRequestBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, putResourcesIdMeasuresBadRequestBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *PutResourcesIDMeasuresBadRequestBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresBadRequest"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("putResourcesIdMeasuresBadRequest"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}

/*PutResourcesIDMeasuresInternalServerErrorBody put resources ID measures internal server error body
swagger:model PutResourcesIDMeasuresInternalServerErrorBody
*/
type PutResourcesIDMeasuresInternalServerErrorBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this put resources ID measures internal server error body
func (o *PutResourcesIDMeasuresInternalServerErrorBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := o.validateType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutResourcesIDMeasuresInternalServerErrorBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresInternalServerError"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var putResourcesIdMeasuresInternalServerErrorBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		putResourcesIdMeasuresInternalServerErrorBodyTypeTypePropEnum = append(putResourcesIdMeasuresInternalServerErrorBodyTypeTypePropEnum, v)
	}
}

const (
	// PutResourcesIDMeasuresInternalServerErrorBodyTypeBadRequest captures enum value "bad_request"
	PutResourcesIDMeasuresInternalServerErrorBodyTypeBadRequest string = "bad_request"
	// PutResourcesIDMeasuresInternalServerErrorBodyTypeUnauthorized captures enum value "unauthorized"
	PutResourcesIDMeasuresInternalServerErrorBodyTypeUnauthorized string = "unauthorized"
	// PutResourcesIDMeasuresInternalServerErrorBodyTypeNotFound captures enum value "not_found"
	PutResourcesIDMeasuresInternalServerErrorBodyTypeNotFound string = "not_found"
	// PutResourcesIDMeasuresInternalServerErrorBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	PutResourcesIDMeasuresInternalServerErrorBodyTypeMethodNotAllowed string = "method_not_allowed"
	// PutResourcesIDMeasuresInternalServerErrorBodyTypeInternal captures enum value "internal"
	PutResourcesIDMeasuresInternalServerErrorBodyTypeInternal string = "internal"
	// PutResourcesIDMeasuresInternalServerErrorBodyTypeInvalidGrant captures enum value "invalid_grant"
	PutResourcesIDMeasuresInternalServerErrorBodyTypeInvalidGrant string = "invalid_grant"
	// PutResourcesIDMeasuresInternalServerErrorBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	PutResourcesIDMeasuresInternalServerErrorBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *PutResourcesIDMeasuresInternalServerErrorBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, putResourcesIdMeasuresInternalServerErrorBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *PutResourcesIDMeasuresInternalServerErrorBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresInternalServerError"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("putResourcesIdMeasuresInternalServerError"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}

/*PutResourcesIDMeasuresNotFoundBody put resources ID measures not found body
swagger:model PutResourcesIDMeasuresNotFoundBody
*/
type PutResourcesIDMeasuresNotFoundBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this put resources ID measures not found body
func (o *PutResourcesIDMeasuresNotFoundBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := o.validateType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutResourcesIDMeasuresNotFoundBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresNotFound"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var putResourcesIdMeasuresNotFoundBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		putResourcesIdMeasuresNotFoundBodyTypeTypePropEnum = append(putResourcesIdMeasuresNotFoundBodyTypeTypePropEnum, v)
	}
}

const (
	// PutResourcesIDMeasuresNotFoundBodyTypeBadRequest captures enum value "bad_request"
	PutResourcesIDMeasuresNotFoundBodyTypeBadRequest string = "bad_request"
	// PutResourcesIDMeasuresNotFoundBodyTypeUnauthorized captures enum value "unauthorized"
	PutResourcesIDMeasuresNotFoundBodyTypeUnauthorized string = "unauthorized"
	// PutResourcesIDMeasuresNotFoundBodyTypeNotFound captures enum value "not_found"
	PutResourcesIDMeasuresNotFoundBodyTypeNotFound string = "not_found"
	// PutResourcesIDMeasuresNotFoundBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	PutResourcesIDMeasuresNotFoundBodyTypeMethodNotAllowed string = "method_not_allowed"
	// PutResourcesIDMeasuresNotFoundBodyTypeInternal captures enum value "internal"
	PutResourcesIDMeasuresNotFoundBodyTypeInternal string = "internal"
	// PutResourcesIDMeasuresNotFoundBodyTypeInvalidGrant captures enum value "invalid_grant"
	PutResourcesIDMeasuresNotFoundBodyTypeInvalidGrant string = "invalid_grant"
	// PutResourcesIDMeasuresNotFoundBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	PutResourcesIDMeasuresNotFoundBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *PutResourcesIDMeasuresNotFoundBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, putResourcesIdMeasuresNotFoundBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *PutResourcesIDMeasuresNotFoundBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresNotFound"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("putResourcesIdMeasuresNotFound"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}

/*PutResourcesIDMeasuresUnauthorizedBody put resources ID measures unauthorized body
swagger:model PutResourcesIDMeasuresUnauthorizedBody
*/
type PutResourcesIDMeasuresUnauthorizedBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this put resources ID measures unauthorized body
func (o *PutResourcesIDMeasuresUnauthorizedBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := o.validateType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutResourcesIDMeasuresUnauthorizedBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresUnauthorized"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var putResourcesIdMeasuresUnauthorizedBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		putResourcesIdMeasuresUnauthorizedBodyTypeTypePropEnum = append(putResourcesIdMeasuresUnauthorizedBodyTypeTypePropEnum, v)
	}
}

const (
	// PutResourcesIDMeasuresUnauthorizedBodyTypeBadRequest captures enum value "bad_request"
	PutResourcesIDMeasuresUnauthorizedBodyTypeBadRequest string = "bad_request"
	// PutResourcesIDMeasuresUnauthorizedBodyTypeUnauthorized captures enum value "unauthorized"
	PutResourcesIDMeasuresUnauthorizedBodyTypeUnauthorized string = "unauthorized"
	// PutResourcesIDMeasuresUnauthorizedBodyTypeNotFound captures enum value "not_found"
	PutResourcesIDMeasuresUnauthorizedBodyTypeNotFound string = "not_found"
	// PutResourcesIDMeasuresUnauthorizedBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	PutResourcesIDMeasuresUnauthorizedBodyTypeMethodNotAllowed string = "method_not_allowed"
	// PutResourcesIDMeasuresUnauthorizedBodyTypeInternal captures enum value "internal"
	PutResourcesIDMeasuresUnauthorizedBodyTypeInternal string = "internal"
	// PutResourcesIDMeasuresUnauthorizedBodyTypeInvalidGrant captures enum value "invalid_grant"
	PutResourcesIDMeasuresUnauthorizedBodyTypeInvalidGrant string = "invalid_grant"
	// PutResourcesIDMeasuresUnauthorizedBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	PutResourcesIDMeasuresUnauthorizedBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *PutResourcesIDMeasuresUnauthorizedBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, putResourcesIdMeasuresUnauthorizedBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *PutResourcesIDMeasuresUnauthorizedBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("putResourcesIdMeasuresUnauthorized"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("putResourcesIdMeasuresUnauthorized"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}
