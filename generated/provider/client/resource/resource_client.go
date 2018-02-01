package resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new resource API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for resource API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
DeleteResourcesID deprovisions

Manifold calls this endpoint to request a resource be deprovisioned.
When a resource is deprovisioned, all attached credentials are assumed
to be deprovisioned as well.

If the resource has already been deprovisioned, then the provider
should return a 404 response.

A response should only be returned once an error has occurred *or* when
the resource is no longer accessible by the user. If a requested action
could take longer than 60s to complete, a callback *must* be used.

**Request Timeout**

If the request takes longer than 60 seconds, then it is assumed to have
failed. Manifold will retry the request again in the future.

**Callback Timeout**

If a `202 Accepted` response is returned, Manifold will expect the
provider to complete the deprovision flow by calling the callback url
within 24 hours. If the callback is not invoked, Manifold will retry
the request again.

If the deprovision was successful, then a `404 Not Found` response
should be returned to Manifold.

*/
func (a *Client) DeleteResourcesID(params *DeleteResourcesIDParams) (*DeleteResourcesIDAccepted, *DeleteResourcesIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteResourcesIDParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DeleteResourcesID",
		Method:             "DELETE",
		PathPattern:        "/resources/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteResourcesIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, nil, err
	}
	switch value := result.(type) {
	case *DeleteResourcesIDAccepted:
		return value, nil, nil
	case *DeleteResourcesIDNoContent:
		return nil, value, nil
	}
	return nil, nil, nil

}

/*
GetResourcesIDMeasures gets how much a resource has used its features

Manifold will call this endpoint daily to get usage information about a
resource and its features.

The Provider should only need to hold information about the time left
in the current month and the previous month.

*/
func (a *Client) GetResourcesIDMeasures(params *GetResourcesIDMeasuresParams) (*GetResourcesIDMeasuresOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetResourcesIDMeasuresParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetResourcesIDMeasures",
		Method:             "GET",
		PathPattern:        "/resources/{id}/measures",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetResourcesIDMeasuresReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetResourcesIDMeasuresOK), nil

}

/*
PatchResourcesID changes plan

Manifold will call this endpoint to request a change in plan of a
resource (either an upgrade or downgrade). This route must support
being called more than once with the same payload.

The `plan` property is the machine readable name of the plan that the
resource is being resized to. The list of possible values are provided
by the provider and stored in the Manifold Catalog.

A response should only be returned once an error has occurred *or* when
the plan change has been completed. If a requested action could take
longer than 60s to complete, a callback *must* be used.

**Request Timeout**

If the request takes longer than 60 seconds, then it is assumed to have
failed. Manifold will retry the request again in the future.

**Callback Timeout**

If a `202 Accepted` response is returned, Manifold will expect the
provider to complete the plan change flow by calling the callback url
within 24 hours. If the callback is not invoked, Manifold will retry
the request again.

If the resource's plan matches the given plan then a `200 Success` or
`204 No Content` response should be returned.

*/
func (a *Client) PatchResourcesID(params *PatchResourcesIDParams) (*PatchResourcesIDOK, *PatchResourcesIDAccepted, *PatchResourcesIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPatchResourcesIDParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "PatchResourcesID",
		Method:             "PATCH",
		PathPattern:        "/resources/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PatchResourcesIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	switch value := result.(type) {
	case *PatchResourcesIDOK:
		return value, nil, nil, nil
	case *PatchResourcesIDAccepted:
		return nil, value, nil, nil
	case *PatchResourcesIDNoContent:
		return nil, nil, value, nil
	}
	return nil, nil, nil, nil

}

/*
PutResourcesID provisions

Manifold will call this endpoint to request the provisiong of a
resource using the provided identifier. This route must support being
called more than once with the same payload.

The `id` property is the unique identifier Manifold will map to this
resource. Providers should use this value for mapping Manifold
Resources to data inside their systems.

The `product`, `plan`, and `region` properties are machine readable
names for the type of product, its plan, and the region in which the
request resource is to be provisioned. These values map to
configuration stored inside the Manifold Catalog.

A response should only be returned once an error has occurred *or* the
provisioned resource is ready for a user to use. If a requested action
could take longer than 60s to complete, a callback *must* be used.

**Request Timeout**

If the request takes longer than 60 seconds, then it is assumed to have
failed. Manifold will retry the request again in the future.

**Callback Timeout**

If a `202 Accepted` response is returned, Manifold will expect the
provider to complete the provision flow by calling the callback url
within 24 hours. If the callback is not invoked, Manifold will retry
the request again.

If the resource has been provisioned successfully with properties that
match the request, then the provider should return a `201 Created` or
`204 No Content` response. However, if the resource provisiong failed,
a corresponding error should be returned.

*/
func (a *Client) PutResourcesID(params *PutResourcesIDParams) (*PutResourcesIDCreated, *PutResourcesIDAccepted, *PutResourcesIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPutResourcesIDParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "PutResourcesID",
		Method:             "PUT",
		PathPattern:        "/resources/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PutResourcesIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	switch value := result.(type) {
	case *PutResourcesIDCreated:
		return value, nil, nil, nil
	case *PutResourcesIDAccepted:
		return nil, value, nil, nil
	case *PutResourcesIDNoContent:
		return nil, nil, value, nil
	}
	return nil, nil, nil, nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
