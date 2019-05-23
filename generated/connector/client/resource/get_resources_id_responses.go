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

	"github.com/manifoldco/grafton/generated/connector/models"
)

// GetResourcesIDReader is a Reader for the GetResourcesID structure.
type GetResourcesIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetResourcesIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetResourcesIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewGetResourcesIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewGetResourcesIDUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewGetResourcesIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewGetResourcesIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetResourcesIDOK creates a GetResourcesIDOK with default headers values
func NewGetResourcesIDOK() *GetResourcesIDOK {
	return &GetResourcesIDOK{}
}

/*GetResourcesIDOK handles this case with default header values.

A resource.
*/
type GetResourcesIDOK struct {
	Payload *models.Resource
}

func (o *GetResourcesIDOK) Error() string {
	return fmt.Sprintf("[GET /resources/{id}][%d] getResourcesIdOK  %+v", 200, o.Payload)
}

func (o *GetResourcesIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Resource)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetResourcesIDBadRequest creates a GetResourcesIDBadRequest with default headers values
func NewGetResourcesIDBadRequest() *GetResourcesIDBadRequest {
	return &GetResourcesIDBadRequest{}
}

/*GetResourcesIDBadRequest handles this case with default header values.

Request denied due to invalid request body, path, or headers.
*/
type GetResourcesIDBadRequest struct {
	Payload GetResourcesIDBadRequestBody
}

func (o *GetResourcesIDBadRequest) Error() string {
	return fmt.Sprintf("[GET /resources/{id}][%d] getResourcesIdBadRequest  %+v", 400, o.Payload)
}

func (o *GetResourcesIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetResourcesIDUnauthorized creates a GetResourcesIDUnauthorized with default headers values
func NewGetResourcesIDUnauthorized() *GetResourcesIDUnauthorized {
	return &GetResourcesIDUnauthorized{}
}

/*GetResourcesIDUnauthorized handles this case with default header values.

Request denied as the provided credentials are no longer valid.
*/
type GetResourcesIDUnauthorized struct {
	Payload GetResourcesIDUnauthorizedBody
}

func (o *GetResourcesIDUnauthorized) Error() string {
	return fmt.Sprintf("[GET /resources/{id}][%d] getResourcesIdUnauthorized  %+v", 401, o.Payload)
}

func (o *GetResourcesIDUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetResourcesIDNotFound creates a GetResourcesIDNotFound with default headers values
func NewGetResourcesIDNotFound() *GetResourcesIDNotFound {
	return &GetResourcesIDNotFound{}
}

/*GetResourcesIDNotFound handles this case with default header values.

Request denied as the requested resource does not exist.
*/
type GetResourcesIDNotFound struct {
	Payload GetResourcesIDNotFoundBody
}

func (o *GetResourcesIDNotFound) Error() string {
	return fmt.Sprintf("[GET /resources/{id}][%d] getResourcesIdNotFound  %+v", 404, o.Payload)
}

func (o *GetResourcesIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetResourcesIDInternalServerError creates a GetResourcesIDInternalServerError with default headers values
func NewGetResourcesIDInternalServerError() *GetResourcesIDInternalServerError {
	return &GetResourcesIDInternalServerError{}
}

/*GetResourcesIDInternalServerError handles this case with default header values.

Request failed due to an internal server error.
*/
type GetResourcesIDInternalServerError struct {
	Payload GetResourcesIDInternalServerErrorBody
}

func (o *GetResourcesIDInternalServerError) Error() string {
	return fmt.Sprintf("[GET /resources/{id}][%d] getResourcesIdInternalServerError  %+v", 500, o.Payload)
}

func (o *GetResourcesIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetResourcesIDBadRequestBody get resources ID bad request body
swagger:model GetResourcesIDBadRequestBody
*/
type GetResourcesIDBadRequestBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this get resources ID bad request body
func (o *GetResourcesIDBadRequestBody) Validate(formats strfmt.Registry) error {
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

func (o *GetResourcesIDBadRequestBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdBadRequest"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var getResourcesIdBadRequestBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		getResourcesIdBadRequestBodyTypeTypePropEnum = append(getResourcesIdBadRequestBodyTypeTypePropEnum, v)
	}
}

const (
	// GetResourcesIDBadRequestBodyTypeBadRequest captures enum value "bad_request"
	GetResourcesIDBadRequestBodyTypeBadRequest string = "bad_request"
	// GetResourcesIDBadRequestBodyTypeUnauthorized captures enum value "unauthorized"
	GetResourcesIDBadRequestBodyTypeUnauthorized string = "unauthorized"
	// GetResourcesIDBadRequestBodyTypeNotFound captures enum value "not_found"
	GetResourcesIDBadRequestBodyTypeNotFound string = "not_found"
	// GetResourcesIDBadRequestBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	GetResourcesIDBadRequestBodyTypeMethodNotAllowed string = "method_not_allowed"
	// GetResourcesIDBadRequestBodyTypeInternal captures enum value "internal"
	GetResourcesIDBadRequestBodyTypeInternal string = "internal"
	// GetResourcesIDBadRequestBodyTypeInvalidGrant captures enum value "invalid_grant"
	GetResourcesIDBadRequestBodyTypeInvalidGrant string = "invalid_grant"
	// GetResourcesIDBadRequestBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	GetResourcesIDBadRequestBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *GetResourcesIDBadRequestBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, getResourcesIdBadRequestBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *GetResourcesIDBadRequestBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdBadRequest"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("getResourcesIdBadRequest"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}

/*GetResourcesIDInternalServerErrorBody get resources ID internal server error body
swagger:model GetResourcesIDInternalServerErrorBody
*/
type GetResourcesIDInternalServerErrorBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this get resources ID internal server error body
func (o *GetResourcesIDInternalServerErrorBody) Validate(formats strfmt.Registry) error {
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

func (o *GetResourcesIDInternalServerErrorBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdInternalServerError"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var getResourcesIdInternalServerErrorBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		getResourcesIdInternalServerErrorBodyTypeTypePropEnum = append(getResourcesIdInternalServerErrorBodyTypeTypePropEnum, v)
	}
}

const (
	// GetResourcesIDInternalServerErrorBodyTypeBadRequest captures enum value "bad_request"
	GetResourcesIDInternalServerErrorBodyTypeBadRequest string = "bad_request"
	// GetResourcesIDInternalServerErrorBodyTypeUnauthorized captures enum value "unauthorized"
	GetResourcesIDInternalServerErrorBodyTypeUnauthorized string = "unauthorized"
	// GetResourcesIDInternalServerErrorBodyTypeNotFound captures enum value "not_found"
	GetResourcesIDInternalServerErrorBodyTypeNotFound string = "not_found"
	// GetResourcesIDInternalServerErrorBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	GetResourcesIDInternalServerErrorBodyTypeMethodNotAllowed string = "method_not_allowed"
	// GetResourcesIDInternalServerErrorBodyTypeInternal captures enum value "internal"
	GetResourcesIDInternalServerErrorBodyTypeInternal string = "internal"
	// GetResourcesIDInternalServerErrorBodyTypeInvalidGrant captures enum value "invalid_grant"
	GetResourcesIDInternalServerErrorBodyTypeInvalidGrant string = "invalid_grant"
	// GetResourcesIDInternalServerErrorBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	GetResourcesIDInternalServerErrorBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *GetResourcesIDInternalServerErrorBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, getResourcesIdInternalServerErrorBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *GetResourcesIDInternalServerErrorBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdInternalServerError"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("getResourcesIdInternalServerError"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}

/*GetResourcesIDNotFoundBody get resources ID not found body
swagger:model GetResourcesIDNotFoundBody
*/
type GetResourcesIDNotFoundBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this get resources ID not found body
func (o *GetResourcesIDNotFoundBody) Validate(formats strfmt.Registry) error {
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

func (o *GetResourcesIDNotFoundBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdNotFound"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var getResourcesIdNotFoundBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		getResourcesIdNotFoundBodyTypeTypePropEnum = append(getResourcesIdNotFoundBodyTypeTypePropEnum, v)
	}
}

const (
	// GetResourcesIDNotFoundBodyTypeBadRequest captures enum value "bad_request"
	GetResourcesIDNotFoundBodyTypeBadRequest string = "bad_request"
	// GetResourcesIDNotFoundBodyTypeUnauthorized captures enum value "unauthorized"
	GetResourcesIDNotFoundBodyTypeUnauthorized string = "unauthorized"
	// GetResourcesIDNotFoundBodyTypeNotFound captures enum value "not_found"
	GetResourcesIDNotFoundBodyTypeNotFound string = "not_found"
	// GetResourcesIDNotFoundBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	GetResourcesIDNotFoundBodyTypeMethodNotAllowed string = "method_not_allowed"
	// GetResourcesIDNotFoundBodyTypeInternal captures enum value "internal"
	GetResourcesIDNotFoundBodyTypeInternal string = "internal"
	// GetResourcesIDNotFoundBodyTypeInvalidGrant captures enum value "invalid_grant"
	GetResourcesIDNotFoundBodyTypeInvalidGrant string = "invalid_grant"
	// GetResourcesIDNotFoundBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	GetResourcesIDNotFoundBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *GetResourcesIDNotFoundBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, getResourcesIdNotFoundBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *GetResourcesIDNotFoundBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdNotFound"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("getResourcesIdNotFound"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}

/*GetResourcesIDUnauthorizedBody get resources ID unauthorized body
swagger:model GetResourcesIDUnauthorizedBody
*/
type GetResourcesIDUnauthorizedBody struct {

	// Explanation of the errors
	// Required: true
	Message []string `json:"message"`

	// The error type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this get resources ID unauthorized body
func (o *GetResourcesIDUnauthorizedBody) Validate(formats strfmt.Registry) error {
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

func (o *GetResourcesIDUnauthorizedBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdUnauthorized"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

var getResourcesIdUnauthorizedBodyTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["bad_request","unauthorized","not_found","method_not_allowed","internal","invalid_grant","unsupported_grant_type"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		getResourcesIdUnauthorizedBodyTypeTypePropEnum = append(getResourcesIdUnauthorizedBodyTypeTypePropEnum, v)
	}
}

const (
	// GetResourcesIDUnauthorizedBodyTypeBadRequest captures enum value "bad_request"
	GetResourcesIDUnauthorizedBodyTypeBadRequest string = "bad_request"
	// GetResourcesIDUnauthorizedBodyTypeUnauthorized captures enum value "unauthorized"
	GetResourcesIDUnauthorizedBodyTypeUnauthorized string = "unauthorized"
	// GetResourcesIDUnauthorizedBodyTypeNotFound captures enum value "not_found"
	GetResourcesIDUnauthorizedBodyTypeNotFound string = "not_found"
	// GetResourcesIDUnauthorizedBodyTypeMethodNotAllowed captures enum value "method_not_allowed"
	GetResourcesIDUnauthorizedBodyTypeMethodNotAllowed string = "method_not_allowed"
	// GetResourcesIDUnauthorizedBodyTypeInternal captures enum value "internal"
	GetResourcesIDUnauthorizedBodyTypeInternal string = "internal"
	// GetResourcesIDUnauthorizedBodyTypeInvalidGrant captures enum value "invalid_grant"
	GetResourcesIDUnauthorizedBodyTypeInvalidGrant string = "invalid_grant"
	// GetResourcesIDUnauthorizedBodyTypeUnsupportedGrantType captures enum value "unsupported_grant_type"
	GetResourcesIDUnauthorizedBodyTypeUnsupportedGrantType string = "unsupported_grant_type"
)

// prop value enum
func (o *GetResourcesIDUnauthorizedBody) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, getResourcesIdUnauthorizedBodyTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *GetResourcesIDUnauthorizedBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("getResourcesIdUnauthorized"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	// value enum
	if err := o.validateTypeEnum("getResourcesIdUnauthorized"+"."+"type", "body", *o.Type); err != nil {
		return err
	}

	return nil
}
