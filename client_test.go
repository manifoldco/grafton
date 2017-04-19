package grafton

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-signature"
)

type stubSigner struct{}

func (stubSigner) Sign([]byte) (*signature.Signature, error) { return &signature.Signature{}, nil }

func callProvision(rawURL string) (string, bool, error) {
	ctx := context.Background()
	sURL, _ := url.Parse(rawURL)
	c := New(sURL, &url.URL{}, stubSigner{})

	cbID := manifold.ID{}
	resID := manifold.ID{}
	return c.ProvisionResource(ctx, cbID, resID, "my-product", "my-plan", "aws::us-east-1")
}

func callProvisionCredentials(rawURL string) (map[string]string, string, bool, error) {
	ctx := context.Background()
	sURL, _ := url.Parse(rawURL)
	c := New(sURL, &url.URL{}, stubSigner{})

	cbID := manifold.ID{}
	resID := manifold.ID{}
	credID := manifold.ID{}
	return c.ProvisionCredentials(ctx, cbID, resID, credID)
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

		gm.Expect(err).To(gm.HaveOccurred())
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

		gm.Expect(err).ToNot(gm.HaveOccurred())
		gm.Expect(async).To(gm.BeFalse(), "Result on 201 should not be async")
		gm.Expect(message).To(gm.BeEmpty(), "No message was expected")
		gm.Expect(creds).To(gm.Equal(map[string]string{"foo": "bar"}))
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

		_, _, _, err := callProvisionCredentials(srv.URL)

		gm.Expect(err).To(gm.HaveOccurred())
	})
}
