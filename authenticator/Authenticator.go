package authenticator

import (
	"net/http"
)

type Authenticator interface {
	// IsAuthenticated returns true iff the request indicates a properly
	// authenticated caller.
	IsAuthenticated(request *http.Request) bool
}
