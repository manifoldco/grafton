package o_auth

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

	"github.com/manifoldco/grafton/generated/connector/models"
)

// NewPostInternalCredentialsParams creates a new PostInternalCredentialsParams object
// with the default values initialized.
func NewPostInternalCredentialsParams() *PostInternalCredentialsParams {
	var ()
	return &PostInternalCredentialsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostInternalCredentialsParamsWithTimeout creates a new PostInternalCredentialsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostInternalCredentialsParamsWithTimeout(timeout time.Duration) *PostInternalCredentialsParams {
	var ()
	return &PostInternalCredentialsParams{

		timeout: timeout,
	}
}

// NewPostInternalCredentialsParamsWithContext creates a new PostInternalCredentialsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostInternalCredentialsParamsWithContext(ctx context.Context) *PostInternalCredentialsParams {
	var ()
	return &PostInternalCredentialsParams{

		Context: ctx,
	}
}

// NewPostInternalCredentialsParamsWithHTTPClient creates a new PostInternalCredentialsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostInternalCredentialsParamsWithHTTPClient(client *http.Client) *PostInternalCredentialsParams {
	var ()
	return &PostInternalCredentialsParams{
		HTTPClient: client,
	}
}

/*PostInternalCredentialsParams contains all the parameters to send to the API endpoint
for the post internal credentials operation typically these are written to a http.Request
*/
type PostInternalCredentialsParams struct {

	/*Body
	  A product id and description for the credential pair.

	*/
	Body *models.OAuthCredentialCreateRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post internal credentials params
func (o *PostInternalCredentialsParams) WithTimeout(timeout time.Duration) *PostInternalCredentialsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post internal credentials params
func (o *PostInternalCredentialsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post internal credentials params
func (o *PostInternalCredentialsParams) WithContext(ctx context.Context) *PostInternalCredentialsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post internal credentials params
func (o *PostInternalCredentialsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post internal credentials params
func (o *PostInternalCredentialsParams) WithHTTPClient(client *http.Client) *PostInternalCredentialsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post internal credentials params
func (o *PostInternalCredentialsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the post internal credentials params
func (o *PostInternalCredentialsParams) WithBody(body *models.OAuthCredentialCreateRequest) *PostInternalCredentialsParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the post internal credentials params
func (o *PostInternalCredentialsParams) SetBody(body *models.OAuthCredentialCreateRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *PostInternalCredentialsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body == nil {
		o.Body = new(models.OAuthCredentialCreateRequest)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
