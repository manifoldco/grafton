package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/validate"

	manifold "github.com/manifoldco/go-manifold"
)

// AccessTokenRequest HTTP Request Body of an Access Token
// swagger:discriminator AccessTokenRequest grant_type
type AccessTokenRequest interface {
	runtime.Validatable

	// client id
	ClientID() manifold.ID
	SetClientID(manifold.ID)

	// client secret
	ClientSecret() OAuthClientSecret
	SetClientSecret(OAuthClientSecret)

	// grant type
	// Required: true
	GrantType() string
	SetGrantType(string)
}

// UnmarshalAccessTokenRequestSlice unmarshals polymorphic slices of AccessTokenRequest
func UnmarshalAccessTokenRequestSlice(reader io.Reader, consumer runtime.Consumer) ([]AccessTokenRequest, error) {
	var elements []json.RawMessage
	if err := consumer.Consume(reader, &elements); err != nil {
		return nil, err
	}

	var result []AccessTokenRequest
	for _, element := range elements {
		obj, err := unmarshalAccessTokenRequest(element, consumer)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

// UnmarshalAccessTokenRequest unmarshals polymorphic AccessTokenRequest
func UnmarshalAccessTokenRequest(reader io.Reader, consumer runtime.Consumer) (AccessTokenRequest, error) {
	// we need to read this twice, so first into a buffer
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return unmarshalAccessTokenRequest(data, consumer)
}

func unmarshalAccessTokenRequest(data []byte, consumer runtime.Consumer) (AccessTokenRequest, error) {
	buf := bytes.NewBuffer(data)
	buf2 := bytes.NewBuffer(data)

	// the first time this is read is to fetch the value of the grant_type property.
	var getType struct {
		GrantType string `json:"grant_type"`
	}
	if err := consumer.Consume(buf, &getType); err != nil {
		return nil, err
	}

	if err := validate.RequiredString("grant_type", "body", getType.GrantType); err != nil {
		return nil, err
	}

	// The value of grant_type is used to determine which type to create and unmarshal the data into
	switch getType.GrantType {
	case "authorization_code":
		var result AuthorizationCode
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case "client_credentials":
		var result ClientCredentials
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil

	}
	return nil, errors.New(422, "invalid grant_type value: %q", getType.GrantType)

}
