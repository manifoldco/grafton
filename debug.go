package grafton

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

// debugRoundTripper implements http.RoundTripper, dumping both request and
// response to stdout
type debugRoundTripper struct {
	rt     http.RoundTripper
	logger *log.Logger
}

// newDebugRoundTripper returns an http.RoundTripper that will dump requests
// and responses before passing them on to the given RoundTripper.
func newDebugRoundTripper(rt http.RoundTripper) *debugRoundTripper {
	return &debugRoundTripper{
		rt:     rt,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (rt *debugRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	dreq, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	rt.logger.Println(string(dreq))

	res, err := rt.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	dres, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}
	rt.logger.Println(string(dres))

	return res, nil
}
