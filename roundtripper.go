package grafton

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/manifoldco/go-signature"
)

// Signer is the interface used to sign outgoing request payloads.
type Signer interface {
	Sign([]byte) (*signature.Signature, error)
}

// SigningRoundTripper implements http.RoundTripper, signing requests before
// sending them.
type signingRoundTripper struct {
	rt     http.RoundTripper
	signer Signer
}

// NewSigningRoundTripper returns an http.RoundTripper that will sign requests
// with the given signer before passing them on to the given RoundTripper.
func newSigningRoundTripper(rt http.RoundTripper, signer Signer) *signingRoundTripper {
	return &signingRoundTripper{
		rt:     rt,
		signer: signer,
	}
}

// RoundTrip implements the http.RoundTripper interface
func (rt *signingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Date", time.Now().UTC().Format(time.RFC3339))
	headers := []string{"host", "date"}

	if req.Header.Get("X-Callback-ID") != "" {
		headers = append(headers, "x-callback-id")
	}

	if req.Header.Get("X-Callback-URL") != "" {
		headers = append(headers, "x-callback-url")
	}

	b := &bytes.Buffer{}
	if req.Body != nil {
		n, err := b.ReadFrom(req.Body)
		if err != nil {
			return nil, err
		}

		defer req.Body.Close()

		if n != 0 {
			// Since we've read from the Body (which is a ReadCloser),
			// we have to replace with a new ReadCloser!
			req.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Length", fmt.Sprintf("%d", b.Len()))
			headers = append(headers, "content-type", "content-length")
		}
	}

	req.Header.Set("X-Signed-Headers", strings.ToLower(strings.Join(headers, " ")))

	canonical, err := signature.Canonize(req, b)
	if err != nil {
		return nil, err
	}

	sig, err := rt.signer.Sign(canonical)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Signature", sig.String())

	return rt.rt.RoundTrip(req)
}
