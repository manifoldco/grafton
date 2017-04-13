// Package grafton provides A simple interface to the provider api.
// it is used both within the grafton test tool, and by our own internal
// services.
package grafton

import (
	"context"
	"errors"
	nurl "net/url"
	"path/filepath"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/manifoldco/go-manifold"

	"github.com/manifoldco/grafton/generated/provider/client"
	"github.com/manifoldco/grafton/generated/provider/client/credential"
	"github.com/manifoldco/grafton/generated/provider/client/resource"
	"github.com/manifoldco/grafton/generated/provider/models"
)

var errMissingMsg = errors.New("`message` field was missing from the response")

// Client is a wrapper around the generated provisioning api client, providing
// convenience methods, and signing outgoing requests.
type Client struct {
	url          *nurl.URL
	connectorURL *nurl.URL
	api          *client.Provider
}

// New creates a new Client
func New(url *nurl.URL, connectorURL *nurl.URL, signer Signer) *Client {
	tp := httptransport.New(url.Host, url.Path, []string{url.Scheme})
	tp.Transport = newSigningRoundTripper(tp.Transport, signer)
	api := client.New(tp, strfmt.Default)

	return &Client{
		url:          url,
		api:          api,
		connectorURL: connectorURL,
	}
}

// ProvisionResource makes a resource provisioning call.
func (c *Client) ProvisionResource(ctx context.Context, cbID, resID manifold.ID, product, plan, region string) (string, bool, error) {

	body := models.ResourceRequest{
		ID:      resID,
		Product: manifold.Label(product),
		Plan:    manifold.Label(plan),
		Region:  models.RegionSlug(region),
	}

	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		return "", false, err
	}

	p := resource.NewPutResourcesIDParams().WithBody(&body).WithID(resID.String())
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)
	res, acceptedRes, _, err := c.api.Resource.PutResourcesID(p)

	var msgPtr *string
	callback := false
	switch {
	case res != nil:
		msgPtr = res.Payload.Message
	case acceptedRes != nil:
		callback = true
		msgPtr = acceptedRes.Payload.Message
	}

	var msg string
	switch {
	case err != nil:
	case msgPtr == nil:
		err = errMissingMsg
	default:
		msg = *msgPtr
	}

	return msg, callback, err
}

func deriveCallbackURL(connectorURL *nurl.URL, cbID manifold.ID) (string, error) {
	u, err := nurl.Parse(connectorURL.String())
	if err != nil {
		return "", err
	}

	u.Path = u.Path + "/callbacks/" + cbID.String()
	return u.String(), nil
}

// ProvisionCredentials makes a credential provisioning call.
func (c *Client) ProvisionCredentials(ctx context.Context, cbID, resID, credID manifold.ID) (map[string]string, string, bool, error) {
	body := models.CredentialRequest{
		ID:         credID,
		ResourceID: resID,
	}

	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		return nil, "", false, err
	}

	p := credential.NewPutCredentialsIDParams().WithBody(&body).WithID(credID.String())
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)

	res, accepted, err := c.api.Credential.PutCredentialsID(p)

	var msg string
	var msgPtr *string
	var creds map[string]string
	callback := accepted != nil
	switch {
	case err != nil:
		return creds, msg, callback, err
	case !callback:
		if res.Payload.Message != nil {
			msgPtr = res.Payload.Message
		}

		creds = res.Payload.Credentials
	case callback && accepted.Payload.Message != nil:
		msgPtr = accepted.Payload.Message
	}

	switch {
	case err != nil:
	case msgPtr == nil:
		err = errMissingMsg
	default:
		msg = *msgPtr
	}

	return creds, msg, callback, err
}

// ChangePlan makes a patch call to change the resource's plan.
func (c *Client) ChangePlan(ctx context.Context, cbID, resourceID manifold.ID, newPlan string) (string, bool, error) {
	body := models.ResourcePlanChangeRequest{Plan: manifold.Label(newPlan)}

	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		return "", false, err
	}

	p := resource.NewPatchResourcesIDParams().WithBody(&body).WithID(resourceID.String())
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)

	res, accepted, _, err := c.api.Resource.PatchResourcesID(p)

	var msg string
	callback := accepted != nil
	switch {
	case res != nil && res.Payload.Message != nil:
		msg = *res.Payload.Message
	case accepted != nil && accepted.Payload.Message != nil:
		msg = *accepted.Payload.Message
	}

	return msg, callback, err
}

// DeprovisionCredentials deletes credentials from the remote provider.
//
// A message will be presented if a callback is provided
func (c *Client) DeprovisionCredentials(ctx context.Context, cbID, credentialID manifold.ID) (string, bool, error) {
	msg := ""
	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		return msg, false, err
	}

	p := credential.NewDeleteCredentialsIDParams().WithID(credentialID.String())
	p.SetContext(ctx)
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)

	accepted, _, err := c.api.Credential.DeleteCredentialsID(p)
	callback := accepted != nil
	if callback && accepted.Payload.Message != nil {
		msg = *accepted.Payload.Message
	}

	return msg, callback, err
}

// DeprovisionResource deletes resources from the remote provider.
func (c *Client) DeprovisionResource(ctx context.Context, cbID, resourceID manifold.ID) (string, bool, error) {
	msg := ""
	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		return msg, false, err
	}

	p := resource.NewDeleteResourcesIDParams().WithID(resourceID.String())
	p.SetContext(ctx)
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)

	accepted, _, err := c.api.Resource.DeleteResourcesID(p)
	callback := accepted != nil
	if callback && accepted.Payload.Message != nil {
		msg = *accepted.Payload.Message
	}

	return msg, callback, err
}

// CreateSsoURL Generates and returns a *url.URL to initiate single sign-on against
// the provider for this client.
func (c *Client) CreateSsoURL(code string, resourceID manifold.ID) *nurl.URL {
	url := *c.url

	url.Path = filepath.Join(url.Path, "/sso/")
	q := nurl.Values{}
	q.Set("code", code)
	q.Set("resource_id", resourceID.String())
	url.RawQuery = q.Encode()

	return &url
}
