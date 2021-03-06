package resource

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
)

// NewGetResourcesIDParams creates a new GetResourcesIDParams object
// with the default values initialized.
func NewGetResourcesIDParams() *GetResourcesIDParams {
	var ()
	return &GetResourcesIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetResourcesIDParamsWithTimeout creates a new GetResourcesIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetResourcesIDParamsWithTimeout(timeout time.Duration) *GetResourcesIDParams {
	var ()
	return &GetResourcesIDParams{

		timeout: timeout,
	}
}

// NewGetResourcesIDParamsWithContext creates a new GetResourcesIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetResourcesIDParamsWithContext(ctx context.Context) *GetResourcesIDParams {
	var ()
	return &GetResourcesIDParams{

		Context: ctx,
	}
}

// NewGetResourcesIDParamsWithHTTPClient creates a new GetResourcesIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetResourcesIDParamsWithHTTPClient(client *http.Client) *GetResourcesIDParams {
	var ()
	return &GetResourcesIDParams{
		HTTPClient: client,
	}
}

/*GetResourcesIDParams contains all the parameters to send to the API endpoint
for the get resources ID operation typically these are written to a http.Request
*/
type GetResourcesIDParams struct {

	/*ID
	  ID of a Resource object, stored as a base32 encoded 18 byte identifier.


	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get resources ID params
func (o *GetResourcesIDParams) WithTimeout(timeout time.Duration) *GetResourcesIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get resources ID params
func (o *GetResourcesIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get resources ID params
func (o *GetResourcesIDParams) WithContext(ctx context.Context) *GetResourcesIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get resources ID params
func (o *GetResourcesIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get resources ID params
func (o *GetResourcesIDParams) WithHTTPClient(client *http.Client) *GetResourcesIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get resources ID params
func (o *GetResourcesIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get resources ID params
func (o *GetResourcesIDParams) WithID(id string) *GetResourcesIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get resources ID params
func (o *GetResourcesIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetResourcesIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
