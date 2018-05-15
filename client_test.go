package grafton

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/errors"
	"github.com/manifoldco/go-manifold/idtype"
	"github.com/manifoldco/go-signature"
)

type stubSigner struct{}

func (stubSigner) Sign([]byte) (*signature.Signature, error) { return &signature.Signature{}, nil }

func callProvision(rawURL string) (string, bool, error) {
	ctx := context.Background()
	sURL, _ := url.Parse(rawURL)
	c := New(sURL, &url.URL{}, stubSigner{}, nil)

	cbID := manifold.ID{}
	resID := manifold.ID{}

	model := ResourceBody{
		ID:         resID,
		Product:    "my-product",
		Plan:       "my-plan",
		Region:     "aws::us-east-1",
		ImportCode: "import-code-1111",
		Features: map[string]interface{}{
			"size":         "40 GB",
			"e-mails":      1000,
			"read-replica": true,
		},
	}

	return c.ProvisionResource(ctx, cbID, model)
}

func callProvisionCredentials(rawURL string) (map[string]string, string, bool, error) {
	ctx := context.Background()
	sURL, _ := url.Parse(rawURL)
	c := New(sURL, &url.URL{}, stubSigner{}, nil)

	cbID := manifold.ID{}
	resID := manifold.ID{}
	credID := manifold.ID{}
	return c.ProvisionCredentials(ctx, cbID, resID, credID)
}

func callChangePlan(rawURL string) (string, bool, error) {
	ctx := context.Background()
	sURL, _ := url.Parse(rawURL)
	c := New(sURL, &url.URL{}, stubSigner{}, nil)

	cbID := manifold.ID{}
	resID := manifold.ID{}

	return c.ChangePlan(ctx, cbID, resID, "new-plan", map[string]interface{}{
		"size":         "40 GB",
		"e-mails":      1000,
		"read-replica": true,
	})
}

func callDeprovisionCredentials(rawURL string) (string, bool, error) {
	ctx := context.Background()
	sURL, _ := url.Parse(rawURL)

	c := New(sURL, &url.URL{}, stubSigner{}, nil)

	cbID := manifold.ID{}
	credID := manifold.ID{}

	return c.DeprovisionCredentials(ctx, cbID, credID)
}

func callDeprovisionResource(rawURL string) (string, bool, error) {
	ctx := context.Background()
	sURL, _ := url.Parse(rawURL)

	c := New(sURL, &url.URL{}, stubSigner{}, nil)

	cbID := manifold.ID{}
	resID := manifold.ID{}

	return c.DeprovisionResource(ctx, cbID, resID)
}

func testCallbackURL(base string) (string, manifold.ID, error) {
	ID, err := manifold.NewID(idtype.Callback)
	if err != nil {
		return "", ID, err
	}

	b, err := url.Parse(base)
	if err != nil {
		return "", ID, err
	}

	c, err := deriveCallbackURL(b, ID)
	if err != nil {
		return "", ID, err
	}

	return c, ID, err
}

func testCreateSSOURL(base, code string) (string, manifold.ID, error) {
	ID, err := manifold.NewID(idtype.OAuthAuthorizationCode)
	if err != nil {
		return "", ID, err
	}

	u, err := url.Parse(base)
	if err != nil {
		return "", ID, err
	}

	c := New(u, &url.URL{}, stubSigner{}, nil)
	return c.CreateSsoURL(code, ID).String(), ID, nil
}

func withCode(code int, fn func(string)) func(t *testing.T) {
	return func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(code)
			rw.Write([]byte(`{"message":"i dont get ya"}`))
		}))
		defer srv.Close()
		fn(srv.URL)
	}
}

func TestDerivingCallbackURL(t *testing.T) {
	t.Run("base url with trailing slash", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, ID, err := testCallbackURL("http://my.host.com/v1/")

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(u).To(gm.Equal("http://my.host.com/v1/callbacks/" + ID.String()))
	})

	t.Run("base url without trailing slash", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, ID, err := testCallbackURL("http://my.host.com/v1")

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(u).To(gm.Equal("http://my.host.com/v1/callbacks/" + ID.String()))
	})
}

func TestDerivingSSOURL(t *testing.T) {
	t.Run("base url with trailing slash", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, ID, err := testCreateSSOURL("http://my.sso.url/v1/", "sdfsd")

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(u).To(gm.Equal("http://my.sso.url/v1/sso?code=sdfsd&resource_id=" + ID.String()))
	})

	t.Run("base url with no trailing slash", func(t *testing.T) {
		gm.RegisterTestingT(t)

		u, ID, err := testCreateSSOURL("http://my.sso.url/v1", "sdfsd")

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(u).To(gm.Equal("http://my.sso.url/v1/sso?code=sdfsd&resource_id=" + ID.String()))
	})
}

func TestProvisionResource(t *testing.T) {
	t.Run("204 no content", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		message, async, err := callProvision(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse(), "Result on 204 should not be async")
		gm.Expect(message).To(gm.BeEmpty(), "No message is returned on 204")
	})

	t.Run("201 with message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(`{"message":"good job"}`))
		}))
		defer srv.Close()

		message, async, err := callProvision(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse(), "Result on 201 should not be async")
		gm.Expect(message).To(gm.Equal("good job"))
	})

	t.Run("201 no message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		_, _, err := callProvision(srv.URL)

		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("202 accepted", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{"message":"please wait"}`))
		}))
		defer srv.Close()

		message, async, err := callProvision(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeTrue(), "Result on 202 should be async")
		gm.Expect(message).To(gm.Equal("please wait"))
	})

	t.Run("204", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		message, async, err := callProvision(srv.URL)
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse())
		gm.Expect(message).To(gm.Equal(""))
	})

	t.Run("400 bad request valid response", withCode(http.StatusBadRequest, func(url string) {
		msg, _, err := callProvision(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.BadRequestError, "i dont get ya")))
	}))

	t.Run("400 bad request with invalid content-type", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "text/html")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"message":"i dont get ya"}`))
		}))
		defer srv.Close()

		msg, _, err := callProvision(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(`no consumer: "text/html"`))
	})

	t.Run("400 bad request with missing message", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		msg, _, err := callProvision(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("401 unauthorized valid response", withCode(http.StatusUnauthorized, func(url string) {
		msg, _, err := callProvision(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.UnauthorizedError, "i dont get ya")))
	}))

	t.Run("409 conflict valid response", withCode(http.StatusConflict, func(url string) {
		msg, _, err := callProvision(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.ConflictError, "i dont get ya")))
	}))

	t.Run("500 internal server error valid response", withCode(http.StatusInternalServerError, func(url string) {
		msg, _, err := callProvision(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.InternalServerError, "i dont get ya")))
	}))

	t.Run("500 internal server error no body", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		msg, async, err := callProvision(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
		gm.Expect(async).To(gm.BeFalse())
	})

	t.Run("503 service unavailable, unrecognized status code", withCode(http.StatusServiceUnavailable, func(url string) {
		msg, _, err := callProvision(url)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
	}))

	t.Run("503 service unavailable, unrecognized status code, bad content-type", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "text/html")
			rw.WriteHeader(http.StatusServiceUnavailable)
			rw.Write([]byte(`{"message":"i dont get ya"}`))
		}))
		defer srv.Close()

		msg, _, err := callProvision(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
	})
}

func TestProvisionCredentials(t *testing.T) {
	t.Run("201 no message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(`{
				"credentials": {
					"foo": "bar"
				}
			}`))
		}))
		defer srv.Close()

		creds, message, async, err := callProvisionCredentials(srv.URL)

		gm.Expect(message).To(gm.Equal(""))
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(creds).To(gm.Equal(map[string]string{"foo": "bar"}))
		gm.Expect(async).To(gm.BeFalse(), "Result on 201 should not be async")
	})

	t.Run("201 with message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(`{
				"message": "have some credentials",
				"credentials": {
					"foo": "bar"
				}
			}`))

		}))
		defer srv.Close()

		creds, message, async, err := callProvisionCredentials(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse(), "Result on 201 should not be async")
		gm.Expect(message).To(gm.Equal("have some credentials"))
		gm.Expect(creds).To(gm.Equal(map[string]string{"foo": "bar"}))
	})

	t.Run("202 accepted", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{"message":"please wait"}`))
		}))
		defer srv.Close()

		creds, message, async, err := callProvisionCredentials(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeTrue(), "Result on 202 should be async")
		gm.Expect(message).To(gm.Equal("please wait"))
		gm.Expect(creds).To(gm.BeEmpty(), "No creds should be returned on 202")
	})

	t.Run("202 no message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		_, msg, _, err := callProvisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("400 bad request valid response", withCode(http.StatusBadRequest, func(url string) {
		_, msg, _, err := callProvisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.BadRequestError, "i dont get ya")))
	}))

	t.Run("400 bad request with invalid content-type", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "text/html")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"message":"i dont get ya"}`))
		}))
		defer srv.Close()

		_, msg, _, err := callProvisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(`no consumer: "text/html"`))
	})

	t.Run("400 bad request with missing message", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		_, msg, _, err := callProvisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("404 not found valid response", withCode(http.StatusNotFound, func(url string) {
		_, msg, _, err := callProvisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.NotFoundError, "i dont get ya")))
	}))

	t.Run("409 conflict valid response", withCode(http.StatusConflict, func(url string) {
		_, msg, _, err := callProvisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.ConflictError, "i dont get ya")))
	}))

	t.Run("500 internal server error valid response", withCode(http.StatusInternalServerError, func(url string) {
		_, msg, _, err := callProvisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.InternalServerError, "i dont get ya")))
	}))

	t.Run("500 internal server error no body", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		_, msg, async, err := callProvisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
		gm.Expect(async).To(gm.BeFalse())
	})
}

func TestChangePlan(t *testing.T) {
	t.Run("200 no message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		message, async, err := callChangePlan(srv.URL)

		gm.Expect(message).To(gm.BeEmpty())
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
		gm.Expect(async).To(gm.BeFalse(), "Result on 200 should not be async")
	})

	t.Run("200 with message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`{
				"message": "the plan has changed"
			}`))

		}))
		defer srv.Close()

		message, async, err := callChangePlan(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse(), "Result on 200 should not be async")
		gm.Expect(message).To(gm.Equal("the plan has changed"))
	})

	t.Run("202 accepted", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{"message":"please wait"}`))
		}))
		defer srv.Close()

		message, async, err := callChangePlan(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeTrue(), "Result on 202 should be async")
		gm.Expect(message).To(gm.Equal("please wait"))
	})

	t.Run("202 no message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		msg, async, err := callChangePlan(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(async).To(gm.BeFalse())
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("204", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		message, async, err := callChangePlan(srv.URL)
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse())
		gm.Expect(message).To(gm.Equal(""))
	})

	t.Run("400 bad request valid response", withCode(http.StatusBadRequest, func(url string) {
		msg, _, err := callChangePlan(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.BadRequestError, "i dont get ya")))
	}))

	t.Run("400 bad request with invalid content-type", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "text/html")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"message":"i dont get ya"}`))
		}))
		defer srv.Close()

		msg, _, err := callChangePlan(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(`no consumer: "text/html"`))
	})

	t.Run("400 bad request with missing message", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		msg, _, err := callChangePlan(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("401 bad unauthorized valid response", withCode(http.StatusUnauthorized, func(url string) {
		msg, _, err := callChangePlan(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.UnauthorizedError, "i dont get ya")))
	}))

	t.Run("404 not found valid response", withCode(http.StatusNotFound, func(url string) {
		msg, _, err := callChangePlan(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.NotFoundError, "i dont get ya")))
	}))

	t.Run("500 internal server error valid response", withCode(http.StatusInternalServerError, func(url string) {
		msg, _, err := callChangePlan(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.InternalServerError, "i dont get ya")))
	}))

	t.Run("500 internal server error no body", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		msg, async, err := callChangePlan(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
		gm.Expect(async).To(gm.BeFalse())
	})
}

func TestDeprovisionCredentials(t *testing.T) {
	t.Run("202 accepted", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{"message":"please wait"}`))
		}))
		defer srv.Close()

		message, async, err := callDeprovisionCredentials(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeTrue(), "Result on 202 should be async")
		gm.Expect(message).To(gm.Equal("please wait"))
	})

	t.Run("202 no message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		msg, async, err := callDeprovisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
		gm.Expect(async).To(gm.BeFalse())
	})

	t.Run("204", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		message, async, err := callDeprovisionCredentials(srv.URL)
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse())
		gm.Expect(message).To(gm.Equal(""))
	})

	t.Run("400 bad request error valid response", withCode(http.StatusBadRequest, func(url string) {
		msg, _, err := callDeprovisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.BadRequestError, "i dont get ya")))
	}))

	t.Run("400 bad request with invalid content-type", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "text/html")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"message":"i dont get ya"}`))
		}))
		defer srv.Close()

		msg, _, err := callDeprovisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(`no consumer: "text/html"`))
	})

	t.Run("400 bad request with missing message", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		msg, _, err := callDeprovisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("401 unauthorized valid response", withCode(http.StatusUnauthorized, func(url string) {
		msg, _, err := callDeprovisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.UnauthorizedError, "i dont get ya")))
	}))

	t.Run("404 not found valid response", withCode(http.StatusNotFound, func(url string) {
		msg, _, err := callDeprovisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.NotFoundError, "i dont get ya")))
	}))

	t.Run("500 internal server error valid response", withCode(http.StatusInternalServerError, func(url string) {
		msg, _, err := callDeprovisionCredentials(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.InternalServerError, "i dont get ya")))
	}))

	t.Run("500 internal server error no body", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		msg, async, err := callDeprovisionCredentials(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
		gm.Expect(async).To(gm.BeFalse())
	})
}

func TestDeprovisionResource(t *testing.T) {
	t.Run("202 accepted", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{"message":"please wait"}`))
		}))
		defer srv.Close()

		message, async, err := callDeprovisionResource(srv.URL)

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeTrue(), "Result on 202 should be async")
		gm.Expect(message).To(gm.Equal("please wait"))
	})

	t.Run("202 no message", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		msg, async, err := callDeprovisionResource(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
		gm.Expect(async).To(gm.BeFalse())
	})

	t.Run("204", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		message, async, err := callDeprovisionResource(srv.URL)
		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse())
		gm.Expect(message).To(gm.Equal(""))
	})

	t.Run("400 bad request valid response", withCode(http.StatusBadRequest, func(url string) {
		msg, _, err := callDeprovisionResource(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.BadRequestError, "i dont get ya")))
	}))

	t.Run("400 bad request with invalid content-type", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "text/html")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"message":"i dont get ya"}`))
		}))
		defer srv.Close()

		msg, _, err := callDeprovisionResource(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(`no consumer: "text/html"`))
	})

	t.Run("400 bad request with missing message", func(t *testing.T) {
		gm.RegisterTestingT(t)
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{}`))
		}))
		defer srv.Close()

		msg, _, err := callDeprovisionResource(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(err).To(gm.MatchError(ErrMissingMsg))
	})

	t.Run("400 unauthorized valid response", withCode(http.StatusUnauthorized, func(url string) {
		msg, _, err := callDeprovisionResource(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.UnauthorizedError, "i dont get ya")))
	}))

	t.Run("404 not found valid response", withCode(http.StatusNotFound, func(url string) {
		msg, _, err := callDeprovisionResource(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.NotFoundError, "i dont get ya")))
	}))

	t.Run("500 internal server error valid response", withCode(http.StatusInternalServerError, func(url string) {
		msg, _, err := callDeprovisionResource(url)

		gm.Expect(msg).To(gm.Equal("i dont get ya"))
		gm.Expect(err).To(gm.MatchError(NewError(errors.InternalServerError, "i dont get ya")))
	}))

	t.Run("500 internal server error no body", func(t *testing.T) {
		gm.RegisterTestingT(t)

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		msg, async, err := callDeprovisionResource(srv.URL)

		gm.Expect(msg).To(gm.Equal(""))
		gm.Expect(IsFatal(err)).To(gm.BeFalse())
		gm.Expect(async).To(gm.BeFalse())
	})
}
