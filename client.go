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

// ResourceBody is an exported type that enables external users to pass data
// to the ProvisionResource function as a model.
type ResourceBody struct {
	ID         manifold.ID
	Product    string
	Plan       string
	Region     string
	ImportCode string
	Features   map[string]interface{}
	PlatformID *manifold.ID
}

// New creates a new Client for Grafton.
// Deprecated in favor of NewClient.
func New(url *nurl.URL, connectorURL *nurl.URL, signer Signer, log *logrus.Entry) *Client {
	opt := ClientOptions{
		URL:          url,
		ConnectorURL: connectorURL,
		Signer:       signer,
		Log:          log,
	}

	return NewClient(opt)
}

// ClientOptions is the options to configure Grafton client.
type ClientOptions struct {
	URL          *nurl.URL
	ConnectorURL *nurl.URL
	Debug        bool
	Signer       Signer
	Log          *logrus.Entry
}

// NewClient creates a new Client for Grafton.
func NewClient(opt ClientOptions) *Client {
	tp := httptransport.New(opt.URL.Host, opt.URL.Path, []string{opt.URL.Scheme})

	if opt.Debug {
		debug := newDebugRoundTripper(tp.Transport)
		signing := newSigningRoundTripper(debug, opt.Signer)
		tp.Transport = signing
	} else {
		signing := newSigningRoundTripper(tp.Transport, opt.Signer)
		tp.Transport = signing
	}

	api := client.New(tp, strfmt.Default)

	if opt.Log == nil {
		opt.Log = logrus.NewEntry(nullLogger)
	}

	return &Client{
		url:          opt.URL,
		api:          api,
		connectorURL: opt.ConnectorURL,
		log:          opt.Log,
	}
}

// ProvisionResource makes a resource provisioning call.
//
// A message will be returned if a callback was used *or* a provider returned
// an error with an explanation.
func (c *Client) ProvisionResource(ctx context.Context, cbID manifold.ID,
	model ResourceBody) (string, bool, error) {

	body := models.ResourceRequest{
		ID:         model.ID,
		Product:    manifold.Label(model.Product),
		Plan:       manifold.Label(model.Plan),
		Region:     models.RegionSlug(model.Region),
		Features:   model.Features,
		PlatformID: model.PlatformID,
	}

	cbURL, err := deriveCallbackURL(c.connectorURL, cbID)
	if err != nil {
		c.log.WithError(err).Error("Could not derive callback url")
		return "", false, err
	}

	p := resource.NewPutResourcesIDParams().WithBody(&body).WithID(model.ID.String())
	// No need to set the Date, this is done in roundtripper.go
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)

	c.log.WithFields(logrus.Fields{
		"url":         c.url,
		"resource_id": model.ID,
		"platform_id": model.PlatformID,
		"product":     model.Product,
		"plan":        model.Plan,
		"region":      model.Region,
	}).Info("Sending PUT resource/{id} request to provider")

	res, acceptedRes, noContent, err := c.api.Resource.PutResourcesID(p)

	if err != nil {
		var graftonErr error
		statusCode := 0
		switch e := err.(type) {
		case *resource.PutResourcesIDBadRequest:
			statusCode = 400
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.PutResourcesIDUnauthorized:
			statusCode = 401
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.PutResourcesIDConflict:
			statusCode = 402
			graftonErr = NewErrWithMsg(merrors.ConflictError, e.Payload.Message)
		case *resource.PutResourcesIDInternalServerError:
			statusCode = 500
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Info("Error unrecognized, returning directly")
			return "", false, err
		}

		c.log.WithError(graftonErr).WithField("status_code", statusCode).Error("Received an error from provider")
		return graftonErr.Error(), false, graftonErr
	}

	var msgPtr *string
	callback := false
	switch {
	case res != nil:
		c.log.WithField("status_code", 201).Info("Received response from provider")
		msgPtr = res.Payload.Message
	case acceptedRes != nil:
		c.log.WithFields(logrus.Fields{
			"status_code":  202,
			"callback_url": cbURL,
		}).Info("Received response from provider, will be awaiting a callback")
		callback = true
		msgPtr = acceptedRes.Payload.Message
	case noContent != nil:
		c.log.WithField("status_code", 204).Info("Received response from provider, no content")
		return "", false, nil
	}

	if msgPtr == nil {
		c.log.Error("Received no message from the provider and the response was not a 204. Failing due to missing message.")
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
	// No need to set the Date, this is done in roundtripper.go
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)

	c.log.WithFields(logrus.Fields{
		"url":           c.url,
		"resource_id":   resID,
		"credential_id": credID,
	}).Info("Sending PUT credentials/{id} request to provider")

	res, accepted, err := c.api.Credential.PutCredentialsID(p)

	if err != nil {
		var graftonErr error
		statusCode := 0
		switch e := err.(type) {
		case *credential.PutCredentialsIDBadRequest:
			statusCode = 400
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *credential.PutCredentialsIDUnauthorized:
			statusCode = 401
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *credential.PutCredentialsIDConflict:
			statusCode = 402
			graftonErr = NewErrWithMsg(merrors.ConflictError, e.Payload.Message)
		case *credential.PutCredentialsIDNotFound:
			statusCode = 404
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *credential.PutCredentialsIDInternalServerError:
			statusCode = 500
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Error unrecognized, returning directly")
			return nil, "", false, err
		}

		c.log.WithError(graftonErr).WithField("status_code", statusCode).Error("Received an error from provider")
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

		c.log.WithField("status_code", 201).Info("Received response from provider")
		creds = res.Payload.Credentials
	case accepted != nil:
		c.log.WithFields(logrus.Fields{
			"status_code":  202,
			"callback_url": cbURL,
		}).Info("Received response from provider, will be awaiting a callback")

		// A message must be provided on a 202 Response
		if accepted.Payload.Message == nil {
			c.log.Error("Received no message from the provider. Failing due to missing message.")
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
	// No need to set the Date, this is done in roundtripper.go
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)
	p.SetContext(ctx)

	c.log.WithFields(logrus.Fields{
		"url":         c.url,
		"resource_id": resourceID,
		"new_plan":    newPlan,
	}).Info("Sending PATCH resource/{id} request to provider")

	res, accepted, noContent, err := c.api.Resource.PatchResourcesID(p)
	if err != nil {
		var graftonErr error
		statusCode := 0
		switch e := err.(type) {
		case *resource.PatchResourcesIDBadRequest:
			statusCode = 400
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.PatchResourcesIDNotFound:
			statusCode = 404
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *resource.PatchResourcesIDUnauthorized:
			statusCode = 401
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.PatchResourcesIDInternalServerError:
			statusCode = 500
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Unrecognized error, returning directly")
			return "", false, err
		}

		c.log.WithError(graftonErr).WithField("status_code", statusCode).Error("Received an error from provider")
		return graftonErr.Error(), false, graftonErr
	}

	var msgPtr *string
	callback := accepted != nil
	switch {
	case res != nil:
		c.log.WithField("status_code", 201).Info("Received response from provider")
		msgPtr = res.Payload.Message
	case accepted != nil:
		c.log.WithFields(logrus.Fields{
			"status_code":  202,
			"callback_url": cbURL,
		}).Info("Received response from provider, will be awaiting a callback")
		msgPtr = accepted.Payload.Message
	case noContent != nil:
		c.log.WithField("status_code", 204).Info("Received response from provider, no content")
		return "", false, nil
	}

	if msgPtr == nil {
		c.log.Error("Received no message from the provider and the response was not a 204. Failing due to missing message.")
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
	// No need to set the Date, this is done in roundtripper.go
	p.SetContext(ctx)
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)

	c.log.WithFields(logrus.Fields{
		"url":           c.url,
		"credential_id": credentialID,
	}).Info("Sending DELETE credentials/{id} request to provider")

	accepted, _, err := c.api.Credential.DeleteCredentialsID(p)
	if err != nil {
		var graftonErr error
		statusCode := 0
		switch e := err.(type) {
		case *credential.DeleteCredentialsIDBadRequest:
			statusCode = 400
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *credential.DeleteCredentialsIDNotFound:
			statusCode = 404
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *credential.DeleteCredentialsIDUnauthorized:
			statusCode = 401
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *credential.DeleteCredentialsIDInternalServerError:
			statusCode = 500
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Unrecognized error, returning directly")
			return "", false, err
		}

		c.log.WithError(graftonErr).WithField("status_code", statusCode).Error("Received an error from provider")
		return graftonErr.Error(), false, graftonErr
	}

	callback := accepted != nil
	if callback {
		c.log.WithFields(logrus.Fields{
			"status_code":  202,
			"callback_url": cbURL,
		}).Info("Received response from provider, will be awaiting a callback")

		if accepted.Payload.Message == nil {
			c.log.Error("Received no message from the provider. Failing due to missing message.")
			return "", false, ErrMissingMsg
		}

		msg = *accepted.Payload.Message
	} else {
		c.log.WithField("status_code", 201).Info("Received response from provider")
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
	// No need to set the Date, this is done in roundtripper.go
	p.SetContext(ctx)
	p.SetXCallbackID(cbID.String())
	p.SetXCallbackURL(cbURL)

	c.log.WithFields(logrus.Fields{
		"url":         c.url,
		"resource_id": resourceID,
	}).Info("Sending DELETE resource/{id} request to provider")

	accepted, _, err := c.api.Resource.DeleteResourcesID(p)
	if err != nil {
		var graftonErr error
		statusCode := 0
		switch e := err.(type) {
		case *resource.DeleteResourcesIDBadRequest:
			statusCode = 400
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.DeleteResourcesIDNotFound:
			statusCode = 404
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *resource.DeleteResourcesIDUnauthorized:
			statusCode = 401
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.DeleteResourcesIDInternalServerError:
			statusCode = 500
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Error("Unrecognized error, returning directly")
			return "", false, err
		}

		c.log.WithError(graftonErr).WithField("status_code", statusCode).Error("Received an error from provider")
		return graftonErr.Error(), false, graftonErr
	}

	callback := accepted != nil
	if callback {
		c.log.WithFields(logrus.Fields{
			"status_code":  202,
			"callback_url": cbURL,
		}).Info("Received response from provider, will be awaiting a callback")

		if accepted.Payload.Message == nil {
			c.log.Error("Received no message from the provider. Failing due to missing message.")
			return "", false, ErrMissingMsg
		}

		msg = *accepted.Payload.Message
	} else {
		c.log.WithField("status_code", 201).Info("Received response from provider")
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

	periodStart := strfmt.DateTime(start)
	periodEnd := strfmt.DateTime(end)

	p := resource.NewGetResourcesIDMeasuresParamsWithContext(ctx)
	p.SetID(rid.String())
	p.SetPeriodStart(periodStart)
	p.SetPeriodEnd(periodEnd)

	c.log.WithFields(logrus.Fields{
		"url":          c.url,
		"resource_id":  rid.String(),
		"period_start": start.String(),
		"period_end":   end.String(),
	}).Info("Sending GET resource/{id}/measures request to provider")

	content, empty, err := c.api.Resource.GetResourcesIDMeasures(p)

	if err != nil {
		var graftonErr error
		statusCode := 0
		switch e := err.(type) {
		case *resource.GetResourcesIDMeasuresBadRequest:
			statusCode = 400
			graftonErr = NewErrWithMsg(merrors.BadRequestError, e.Payload.Message)
		case *resource.GetResourcesIDMeasuresUnauthorized:
			statusCode = 401
			graftonErr = NewErrWithMsg(merrors.UnauthorizedError, e.Payload.Message)
		case *resource.GetResourcesIDMeasuresNotFound:
			statusCode = 404
			graftonErr = NewErrWithMsg(merrors.NotFoundError, e.Payload.Message)
		case *resource.GetResourcesIDMeasuresInternalServerError:
			statusCode = 500
			graftonErr = NewErrWithMsg(merrors.InternalServerError, e.Payload.Message)
		default:
			c.log.WithError(err).Info("Error unrecognized, returning directly")
			return nil, err
		}

		c.log.WithError(graftonErr).WithField("status_code", statusCode).Error("Received an error from provider")
		return nil, graftonErr
	}

	if empty != nil {
		c.log.WithField("status_code", 204).Info("Received response from provider")
		return &models.ResourceMeasures{
			ResourceID:  rid,
			PeriodStart: &periodStart,
			PeriodEnd:   &periodEnd,
		}, nil
	}

	c.log.WithField("status_code", 201).Info("Received response from provider")
	return content.Payload, nil
}
