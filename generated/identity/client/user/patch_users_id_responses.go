package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/grafton/generated/identity/models"
)

// PatchUsersIDReader is a Reader for the PatchUsersID structure.
type PatchUsersIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchUsersIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPatchUsersIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewPatchUsersIDBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewPatchUsersIDInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewPatchUsersIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPatchUsersIDOK creates a PatchUsersIDOK with default headers values
func NewPatchUsersIDOK() *PatchUsersIDOK {
	return &PatchUsersIDOK{}
}

/*PatchUsersIDOK handles this case with default header values.

Complete user object
*/
type PatchUsersIDOK struct {
	Payload *models.User
}

func (o *PatchUsersIDOK) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdOK  %+v", 200, o.Payload)
}

func (o *PatchUsersIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.User)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchUsersIDBadRequest creates a PatchUsersIDBadRequest with default headers values
func NewPatchUsersIDBadRequest() *PatchUsersIDBadRequest {
	return &PatchUsersIDBadRequest{}
}

/*PatchUsersIDBadRequest handles this case with default header values.

Validation failed for request
*/
type PatchUsersIDBadRequest struct {
	Payload *manifold.Error
}

func (o *PatchUsersIDBadRequest) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdBadRequest  %+v", 400, o.Payload)
}

func (o *PatchUsersIDBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(manifold.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchUsersIDInternalServerError creates a PatchUsersIDInternalServerError with default headers values
func NewPatchUsersIDInternalServerError() *PatchUsersIDInternalServerError {
	return &PatchUsersIDInternalServerError{}
}

/*PatchUsersIDInternalServerError handles this case with default header values.

Request failed due to an internal server error.
*/
type PatchUsersIDInternalServerError struct {
	Payload *manifold.Error
}

func (o *PatchUsersIDInternalServerError) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] patchUsersIdInternalServerError  %+v", 500, o.Payload)
}

func (o *PatchUsersIDInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(manifold.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchUsersIDDefault creates a PatchUsersIDDefault with default headers values
func NewPatchUsersIDDefault(code int) *PatchUsersIDDefault {
	return &PatchUsersIDDefault{
		_statusCode: code,
	}
}

/*PatchUsersIDDefault handles this case with default header values.

Unexpected error
*/
type PatchUsersIDDefault struct {
	_statusCode int

	Payload *manifold.Error
}

// Code gets the status code for the patch users ID default response
func (o *PatchUsersIDDefault) Code() int {
	return o._statusCode
}

func (o *PatchUsersIDDefault) Error() string {
	return fmt.Sprintf("[PATCH /users/{id}][%d] PatchUsersID default  %+v", o._statusCode, o.Payload)
}

func (o *PatchUsersIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(manifold.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
