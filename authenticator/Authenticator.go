package authenticator

import (
	"net/http"
)

type Authenticator interface {
	// IsAuthenticated returns true iff the request indicates a properly
	// authenticated caller.
	IsAuthenticated(request *http.Request) bool

	// UserId returns a unique identifier of the user being currently logged
	// in.
	// If no user is logged in (therefore, if IsAuthenticated returns false),
	// UserId returns the empty string.
	UserId(request *http.Request) string
}
