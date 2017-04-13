package acceptance

import (
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/onsi/gomega/format"
)

// notError ensures that an error is nil. If its not, it prints out the seen
// error in a nice way.
func notError() *notErrorMatcher {
	return &notErrorMatcher{}
}

type notErrorMatcher struct{}

func (ne notErrorMatcher) Match(actual interface{}) (bool, error) {
	return actual == nil, nil
}

func (ne notErrorMatcher) FailureMessage(actual interface{}) string {
	switch t := actual.(type) {
	case *url.Error:
		return ne.FailureMessage(t.Err)
	case *net.OpError:
		switch t.Err.(type) { //
		case *net.DNSError, *os.SyscallError:
			return t.Error()
		default:
			return ne.FailureMessage(t.Err)
		}
	case *runtime.APIError:
		if t.OperationName == "unknown error" {
			return fmt.Sprintf("unexpected status code '%d' on response", t.Code)
		}

		return t.Error()
	case error:
		return t.Error()
	}

	// Fall back to the default gomega output format
	return format.Message(actual, "To not have occurred.")
}

func (ne notErrorMatcher) NegatedFailureMessage(actual interface{}) string {
	// We don't use the negated form.
	return "Error should have existed"
}
