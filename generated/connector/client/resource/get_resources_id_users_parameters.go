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

// NewGetResourcesIDUsersParams creates a new GetResourcesIDUsersParams object
// with the default values initialized.
func NewGetResourcesIDUsersParams() *GetResourcesIDUsersParams {
	var ()
	return &GetResourcesIDUsersParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetResourcesIDUsersParamsWithTimeout creates a new GetResourcesIDUsersParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetResourcesIDUsersParamsWithTimeout(timeout time.Duration) *GetResourcesIDUsersParams {
	var ()
	return &GetResourcesIDUsersParams{

		timeout: timeout,
	}
}

// NewGetResourcesIDUsersParamsWithContext creates a new GetResourcesIDUsersParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetResourcesIDUsersParamsWithContext(ctx context.Context) *GetResourcesIDUsersParams {
	var ()
	return &GetResourcesIDUsersParams{

		Context: ctx,
	}
}

// NewGetResourcesIDUsersParamsWithHTTPClient creates a new GetResourcesIDUsersParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetResourcesIDUsersParamsWithHTTPClient(client *http.Client) *GetResourcesIDUsersParams {
	var ()
	return &GetResourcesIDUsersParams{
		HTTPClient: client,
	}
}

/*GetResourcesIDUsersParams contains all the parameters to send to the API endpoint
for the get resources ID users operation typically these are written to a http.Request
*/
type GetResourcesIDUsersParams struct {

	/*ID
	  ID of a Resource object, stored as a base32 encoded 18 byte identifier.


	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get resources ID users params
func (o *GetResourcesIDUsersParams) WithTimeout(timeout time.Duration) *GetResourcesIDUsersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get resources ID users params
func (o *GetResourcesIDUsersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get resources ID users params
func (o *GetResourcesIDUsersParams) WithContext(ctx context.Context) *GetResourcesIDUsersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get resources ID users params
func (o *GetResourcesIDUsersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get resources ID users params
func (o *GetResourcesIDUsersParams) WithHTTPClient(client *http.Client) *GetResourcesIDUsersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get resources ID users params
func (o *GetResourcesIDUsersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get resources ID users params
func (o *GetResourcesIDUsersParams) WithID(id string) *GetResourcesIDUsersParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get resources ID users params
func (o *GetResourcesIDUsersParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetResourcesIDUsersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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