package team

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/manifoldco/grafton/generated/identity/models"
)

// NewPatchTeamsIDParams creates a new PatchTeamsIDParams object
// with the default values initialized.
func NewPatchTeamsIDParams() *PatchTeamsIDParams {
	var ()
	return &PatchTeamsIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPatchTeamsIDParamsWithTimeout creates a new PatchTeamsIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPatchTeamsIDParamsWithTimeout(timeout time.Duration) *PatchTeamsIDParams {
	var ()
	return &PatchTeamsIDParams{

		timeout: timeout,
	}
}

// NewPatchTeamsIDParamsWithContext creates a new PatchTeamsIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewPatchTeamsIDParamsWithContext(ctx context.Context) *PatchTeamsIDParams {
	var ()
	return &PatchTeamsIDParams{

		Context: ctx,
	}
}

// NewPatchTeamsIDParamsWithHTTPClient creates a new PatchTeamsIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPatchTeamsIDParamsWithHTTPClient(client *http.Client) *PatchTeamsIDParams {
	var ()
	return &PatchTeamsIDParams{
		HTTPClient: client,
	}
}

/*PatchTeamsIDParams contains all the parameters to send to the API endpoint
for the patch teams ID operation typically these are written to a http.Request
*/
type PatchTeamsIDParams struct {

	/*Body
	  Team update request


	*/
	Body *models.UpdateTeam
	/*ID
	  ID of the Team to lookup, stored as a base32 encoded 18 byte
	identifier.


	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the patch teams ID params
func (o *PatchTeamsIDParams) WithTimeout(timeout time.Duration) *PatchTeamsIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch teams ID params
func (o *PatchTeamsIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch teams ID params
func (o *PatchTeamsIDParams) WithContext(ctx context.Context) *PatchTeamsIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch teams ID params
func (o *PatchTeamsIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch teams ID params
func (o *PatchTeamsIDParams) WithHTTPClient(client *http.Client) *PatchTeamsIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch teams ID params
func (o *PatchTeamsIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the patch teams ID params
func (o *PatchTeamsIDParams) WithBody(body *models.UpdateTeam) *PatchTeamsIDParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the patch teams ID params
func (o *PatchTeamsIDParams) SetBody(body *models.UpdateTeam) {
	o.Body = body
}

// WithID adds the id to the patch teams ID params
func (o *PatchTeamsIDParams) WithID(id string) *PatchTeamsIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the patch teams ID params
func (o *PatchTeamsIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PatchTeamsIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body == nil {
		o.Body = new(models.UpdateTeam)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
