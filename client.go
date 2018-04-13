// Package grafton provides A simple interface to the provider api.
// it is used both within the grafton test tool, and by our own internal
// services.
package grafton

import (
	"context"
	"io/ioutil"
	nurl "net/url"
	"path"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"

	"github.com/manifoldco/go-manifold"
	merrors "github.com/manifoldco/go-manifold/errors"

	"github.com/manifoldco/grafton/generated/provider/client"
	"github.com/manifoldco/grafton/generated/provider/client/credential"
	"github.com/manifoldco/grafton/generated/provider/client/resource"
	"github.com/manifoldco/grafton/generated/provider/models"
)

var nullLogger *logrus.Logger

func init() {
	nullLogger = logrus.New()
	nullLogger.Out = ioutil.Discard
}

// Client is a wrapper around the generated provisioning api client, providing
// convenience methods, and signing outgoing requests.
type Client struct {
	url          *nurl.URL
	connectorURL *nurl.URL
	api          *client.ManifoldProvider
	log          *logrus.Entry
}

// New creates a new Client
func New(url *nurl.URL, connectorURL *nurl.URL, signer Signer, log *logrus.Entry) *Client {
	tp := httptransport.New(url.Host, url.Path, []string{url.Scheme})
	tp.Transport = newSigningRoundTripper(tp.Transport, signer)
	api := client.New(tp, strfmt.Default)

	if log == nil {
		log = logrus.NewEntry(nullLogger)
	}

	return &Client{
		url:          url,
		api:          api,
		connectorURL: connectorURL,
		log:          log,
	}
}

// ProvisionResource makes a resource provisioning call.
//
// A message will be returned if a callback was used *or* a provider returned
// an error with an explanation.
func (c *Client) ProvisionResource(ctx context.Context, cbID, resID manifold.ID, product, plan,
	region string, features map[string]interface{}) (string, bool, error) {

	body := models.ResourceRequest{
		ID:       resID,
		Product:  manifold.Label(product),
		Plan:     manifold.Label(plan),
		Region:   models.RegionSlug(region),
		Features: features,
	}

	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		c.log.WithError(err).Error("Could not derive callback url")
		return "", false, err
	}

	p := resource.NewPutResourcesIDParams().WithBody(&body).WithID(resID.String())
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)
	res, acceptedRes, noContent, err := c.api.Resource.PutResourcesID(p)

	if err != nil {
		var graftonErr error
		switch e := err.(type) {
		case *resource.PutResourcesIDBadRequest:
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.PutResourcesIDUnauthorized:
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.PutResourcesIDConflict:
			graftonErr = NewErrWithMsg(merrors.ConflictError, e.Payload.Message)
		case *resource.PutResourcesIDInternalServerError:
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Info("Error unrecognized, returning directly")
			return "", false, err
		}

		if graftonErr == ErrMissingMsg {
			c.log.WithError(err).Error("Missing error message in response")
			return "", false, graftonErr
		}

		return graftonErr.Error(), false, graftonErr
	}

	var msgPtr *string
	callback := false
	switch {
	case res != nil:
		msgPtr = res.Payload.Message
	case acceptedRes != nil:
		callback = true
		msgPtr = acceptedRes.Payload.Message
	case noContent != nil:
		return "", false, nil
	}

	if msgPtr == nil {
		return "", false, ErrMissingMsg
	}

	return *msgPtr, callback, nil
}

func deriveCallbackURL(connectorURL *nurl.URL, cbID manifold.ID) (string, error) {
	u, err := nurl.Parse(connectorURL.String())
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "callbacks/"+cbID.String())
	return u.String(), nil
}

// ProvisionCredentials makes a credential provisioning call.
//
// A message will be returned if a callback was used *or* a provider returned
// an error with an explanation.
func (c *Client) ProvisionCredentials(ctx context.Context, cbID, resID, credID manifold.ID) (map[string]string, string, bool, error) {
	body := models.CredentialRequest{
		ID:         credID,
		ResourceID: resID,
	}

	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		c.log.WithError(err).Error("Could not derive callback url")
		return nil, "", false, err
	}

	p := credential.NewPutCredentialsIDParams().WithBody(&body).WithID(credID.String())
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)

	res, accepted, err := c.api.Credential.PutCredentialsID(p)

	if err != nil {
		var graftonErr error
		switch e := err.(type) {
		case *credential.PutCredentialsIDBadRequest:
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *credential.PutCredentialsIDUnauthorized:
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *credential.PutCredentialsIDConflict:
			graftonErr = NewErrWithMsg(merrors.ConflictError, e.Payload.Message)
		case *credential.PutCredentialsIDNotFound:
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *credential.PutCredentialsIDInternalServerError:
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Error unrecognized, returning directly")
			return nil, "", false, err
		}

		if graftonErr == ErrMissingMsg {
			c.log.WithError(err).Error("Missing error message in response")
			return nil, "", false, graftonErr
		}

		return nil, graftonErr.Error(), false, graftonErr
	}

	msg := ""
	var creds map[string]string
	callback := accepted != nil
	switch {
	case res != nil:
		// A message is optional on a 201 Response
		if res.Payload.Message != nil {
			msg = *res.Payload.Message
		}

		creds = res.Payload.Credentials
	case accepted != nil:
		// A message must be provided on a 202 Response
		if accepted.Payload.Message == nil {
			return nil, "", false, ErrMissingMsg
		}

		msg = *accepted.Payload.Message
	}

	return creds, msg, callback, err
}

// ChangePlan makes a patch call to change the resource's plan.
//
// A message will be returned if a callback was used *or* a provider returned
// an error with an explanation.
func (c *Client) ChangePlan(ctx context.Context, cbID, resourceID manifold.ID, newPlan string,
	features map[string]interface{}) (string, bool, error) {

	body := models.ResourcePlanChangeRequest{
		Plan:     manifold.Label(newPlan),
		Features: features,
	}

	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		c.log.WithError(err).Error("Error driving callback url")
		return "", false, err
	}

	p := resource.NewPatchResourcesIDParams().WithBody(&body).WithID(resourceID.String())
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)

	res, accepted, noContent, err := c.api.Resource.PatchResourcesID(p)
	if err != nil {
		var graftonErr error
		switch e := err.(type) {
		case *resource.PatchResourcesIDBadRequest:
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.PatchResourcesIDNotFound:
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *resource.PatchResourcesIDUnauthorized:
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.PatchResourcesIDInternalServerError:
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Unrecognized error, returning directly")
			return "", false, err
		}

		if graftonErr == ErrMissingMsg {
			c.log.WithError(err).Error("No error message in response")
			return "", false, graftonErr
		}

		return graftonErr.Error(), false, graftonErr
	}

	var msgPtr *string
	callback := accepted != nil
	switch {
	case res != nil:
		msgPtr = res.Payload.Message
	case accepted != nil:
		msgPtr = accepted.Payload.Message
	case noContent != nil:
		return "", false, nil
	}

	if msgPtr == nil {
		return "", false, ErrMissingMsg
	}

	return *msgPtr, callback, nil
}

// DeprovisionCredentials deletes credentials from the remote provider.
//
// A message will be presented if a callback is provided or if a message was
// returned from the provider due to an error.
func (c *Client) DeprovisionCredentials(ctx context.Context, cbID, credentialID manifold.ID) (string, bool, error) {
	msg := ""
	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		c.log.WithError(err).Error("Could not derive callback url")
		return msg, false, err
	}

	p := credential.NewDeleteCredentialsIDParams().WithID(credentialID.String())
	p.SetContext(ctx)
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)

	accepted, _, err := c.api.Credential.DeleteCredentialsID(p)
	if err != nil {
		var graftonErr error
		switch e := err.(type) {
		case *credential.DeleteCredentialsIDBadRequest:
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *credential.DeleteCredentialsIDNotFound:
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *credential.DeleteCredentialsIDUnauthorized:
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *credential.DeleteCredentialsIDInternalServerError:
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Unrecognized error, returning directly")
			return "", false, err
		}

		if graftonErr == ErrMissingMsg {
			c.log.WithError(err).Error("No error message in response")
			return "", false, graftonErr
		}

		return graftonErr.Error(), false, graftonErr
	}

	callback := accepted != nil
	if callback {
		if accepted.Payload.Message == nil {
			return "", false, ErrMissingMsg
		}

		msg = *accepted.Payload.Message
	}

	return msg, callback, err
}

// DeprovisionResource deletes resources from the remote provider.
//
// A message will be returned if a callback was used *or* a provider returned
// an error with an explanation.
func (c *Client) DeprovisionResource(ctx context.Context, cbID, resourceID manifold.ID) (string, bool, error) {
	msg := ""
	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		c.log.WithError(err).Error("Could not derive callback url")
		return msg, false, err
	}

	p := resource.NewDeleteResourcesIDParams().WithID(resourceID.String())
	p.SetContext(ctx)
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)

	accepted, _, err := c.api.Resource.DeleteResourcesID(p)
	if err != nil {
		var graftonErr error
		switch e := err.(type) {
		case *resource.DeleteResourcesIDBadRequest:
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.DeleteResourcesIDNotFound:
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *resource.DeleteResourcesIDUnauthorized:
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.DeleteResourcesIDInternalServerError:
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Unrecognized error, returning directly")
			return "", false, err
		}

		if graftonErr == ErrMissingMsg {
			c.log.WithError(err).Error("No error message in response")
			return "", false, graftonErr
		}

		return graftonErr.Error(), false, graftonErr
	}

	callback := accepted != nil
	if callback {
		if accepted.Payload.Message == nil {
			return "", false, ErrMissingMsg
		}

		msg = *accepted.Payload.Message
	}

	return msg, callback, err
}

// CreateSsoURL Generates and returns a *url.URL to initiate single sign-on against
// the provider for this client.
func (c *Client) CreateSsoURL(code string, resourceID manifold.ID) *nurl.URL {
	return CreateSsoURL(c.url, code, resourceID)
}

// CreateSsoURL generates and returns a *url.URL to initiate a single sign-on
// request against the provided base url, code, and resourceID.
func CreateSsoURL(base *nurl.URL, code string, resourceID manifold.ID) *nurl.URL {
	url := *base

	url.Path = path.Join(url.Path, "sso/")
	q := nurl.Values{}
	q.Set("code", code)
	q.Set("resource_id", resourceID.String())
	url.RawQuery = q.Encode()

	return &url
}

// PullResourceMeasures tries to get information about a resource usage.
func (c *Client) PullResourceMeasures(ctx context.Context, rid manifold.ID,
	start, end time.Time) (*models.ResourceMeasures, error) {

	p := resource.NewGetResourcesIDMeasuresParamsWithContext(ctx)
	p.SetID(rid.String())
	p.SetPeriodStart(strfmt.DateTime(start))
	p.SetPeriodEnd(strfmt.DateTime(end))

	res, _, err := c.api.Resource.GetResourcesIDMeasures(p)

	if err != nil {
		var graftonErr error
		switch e := err.(type) {
		case *resource.GetResourcesIDMeasuresBadRequest:
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.GetResourcesIDMeasuresUnauthorized:
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.GetResourcesIDMeasuresNotFound:
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *resource.GetResourcesIDMeasuresInternalServerError:
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Info("Error unrecognized, returning directly")
			return nil, err
		}

		if graftonErr == ErrMissingMsg {
			c.log.WithError(err).Error("Missing error message in response")
			return nil, graftonErr
		}

		return nil, graftonErr
	}

	return res.Payload, nil
}
