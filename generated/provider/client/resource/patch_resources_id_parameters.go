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

	"github.com/manifoldco/grafton/generated/provider/models"
)

// NewPatchResourcesIDParams creates a new PatchResourcesIDParams object
// with the default values initialized.
func NewPatchResourcesIDParams() *PatchResourcesIDParams {
	var ()
	return &PatchResourcesIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPatchResourcesIDParamsWithTimeout creates a new PatchResourcesIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPatchResourcesIDParamsWithTimeout(timeout time.Duration) *PatchResourcesIDParams {
	var ()
	return &PatchResourcesIDParams{

		timeout: timeout,
	}
}

// NewPatchResourcesIDParamsWithContext creates a new PatchResourcesIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewPatchResourcesIDParamsWithContext(ctx context.Context) *PatchResourcesIDParams {
	var ()
	return &PatchResourcesIDParams{

		Context: ctx,
	}
}

// NewPatchResourcesIDParamsWithHTTPClient creates a new PatchResourcesIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPatchResourcesIDParamsWithHTTPClient(client *http.Client) *PatchResourcesIDParams {
	var ()
	return &PatchResourcesIDParams{
		HTTPClient: client,
	}
}

/*PatchResourcesIDParams contains all the parameters to send to the API endpoint
for the patch resources ID operation typically these are written to a http.Request
*/
type PatchResourcesIDParams struct {

	/*Date
	  Timestamp of when the request was issued from Manifold in UTC.

	*/
	Date strfmt.DateTime
	/*XCallbackID
	  ID of the Callback for completing this request if the provider returns a
	`202 Accepted`, stored as a base 32 encoded 18 byte identifier.


	*/
	XCallbackID string
	/*XCallbackURL
	  The URL the provider calls to complete the request if a `202 Accepted` is
	returned.


	*/
	XCallbackURL string
	/*XSignature
	  Header containing the signature, ephemeral public key, and
	signature of the used public key signed by the Manifold root
	signing key to prove authenticity of the request.

	```
	X-Signature: [request-signature] [public-key-value] [signature-of-public-key]
	```

	examples:

	```
	X-Signature: 96afMc5FVZjhGQ4YLPyRW9tcYoPKyj1EMUxkzma_Q4c WydEYGQb7Y4ER6KPAc-YuIwAqFUpII5P9U3MAZ3wxAQ ozhcosOmuWltP8r1BAs-0kccV1AkbHcKYLSjU0dGUHY
	```


	*/
	XSignature string
	/*XSignedHeaders
	  Comma delimited ordered list of header fields used in generating
	the request signature.


	*/
	XSignedHeaders string
	/*Body
	  Resource Provisioning Request

	*/
	Body *models.ResourcePlanChangeRequest
	/*ID
	  ID of a Resource object, stored as a base32 encoded 18 byte identifier.


	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the patch resources ID params
func (o *PatchResourcesIDParams) WithTimeout(timeout time.Duration) *PatchResourcesIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch resources ID params
func (o *PatchResourcesIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch resources ID params
func (o *PatchResourcesIDParams) WithContext(ctx context.Context) *PatchResourcesIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch resources ID params
func (o *PatchResourcesIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch resources ID params
func (o *PatchResourcesIDParams) WithHTTPClient(client *http.Client) *PatchResourcesIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch resources ID params
func (o *PatchResourcesIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDate adds the date to the patch resources ID params
func (o *PatchResourcesIDParams) WithDate(date strfmt.DateTime) *PatchResourcesIDParams {
	o.SetDate(date)
	return o
}

// SetDate adds the date to the patch resources ID params
func (o *PatchResourcesIDParams) SetDate(date strfmt.DateTime) {
	o.Date = date
}

// WithXCallbackID adds the xCallbackID to the patch resources ID params
func (o *PatchResourcesIDParams) WithXCallbackID(xCallbackID string) *PatchResourcesIDParams {
	o.SetXCallbackID(xCallbackID)
	return o
}

// SetXCallbackID adds the xCallbackId to the patch resources ID params
func (o *PatchResourcesIDParams) SetXCallbackID(xCallbackID string) {
	o.XCallbackID = xCallbackID
}

// WithXCallbackURL adds the xCallbackURL to the patch resources ID params
func (o *PatchResourcesIDParams) WithXCallbackURL(xCallbackURL string) *PatchResourcesIDParams {
	o.SetXCallbackURL(xCallbackURL)
	return o
}

// SetXCallbackURL adds the xCallbackUrl to the patch resources ID params
func (o *PatchResourcesIDParams) SetXCallbackURL(xCallbackURL string) {
	o.XCallbackURL = xCallbackURL
}

// WithXSignature adds the xSignature to the patch resources ID params
func (o *PatchResourcesIDParams) WithXSignature(xSignature string) *PatchResourcesIDParams {
	o.SetXSignature(xSignature)
	return o
}

// SetXSignature adds the xSignature to the patch resources ID params
func (o *PatchResourcesIDParams) SetXSignature(xSignature string) {
	o.XSignature = xSignature
}

// WithXSignedHeaders adds the xSignedHeaders to the patch resources ID params
func (o *PatchResourcesIDParams) WithXSignedHeaders(xSignedHeaders string) *PatchResourcesIDParams {
	o.SetXSignedHeaders(xSignedHeaders)
	return o
}

// SetXSignedHeaders adds the xSignedHeaders to the patch resources ID params
func (o *PatchResourcesIDParams) SetXSignedHeaders(xSignedHeaders string) {
	o.XSignedHeaders = xSignedHeaders
}

// WithBody adds the body to the patch resources ID params
func (o *PatchResourcesIDParams) WithBody(body *models.ResourcePlanChangeRequest) *PatchResourcesIDParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the patch resources ID params
func (o *PatchResourcesIDParams) SetBody(body *models.ResourcePlanChangeRequest) {
	o.Body = body
}

// WithID adds the id to the patch resources ID params
func (o *PatchResourcesIDParams) WithID(id string) *PatchResourcesIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the patch resources ID params
func (o *PatchResourcesIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PatchResourcesIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param Date
	if err := r.SetHeaderParam("Date", o.Date.String()); err != nil {
		return err
	}

	// header param X-Callback-ID
	if err := r.SetHeaderParam("X-Callback-ID", o.XCallbackID); err != nil {
		return err
	}

	// header param X-Callback-URL
	if err := r.SetHeaderParam("X-Callback-URL", o.XCallbackURL); err != nil {
		return err
	}

	// header param X-Signature
	if err := r.SetHeaderParam("X-Signature", o.XSignature); err != nil {
		return err
	}

	// header param X-Signed-Headers
	if err := r.SetHeaderParam("X-Signed-Headers", o.XSignedHeaders); err != nil {
		return err
	}

	if o.Body == nil {
		o.Body = new(models.ResourcePlanChangeRequest)
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
