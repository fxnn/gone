package authenticator

import "net/http"

// AlwaysAuthenticated is an Authenticator that returns every user as being
// authenticated.
type AlwaysAuthenticated struct{}

// Ensure that Authenticator interface is implemented
var _ Authenticator = (*AlwaysAuthenticated)(nil)

func NewAlwaysAuthenticated() *AlwaysAuthenticated {
	return &AlwaysAuthenticated{}
}

func (AlwaysAuthenticated) IsAuthenticated(request *http.Request) bool {
	return true
}

func (AlwaysAuthenticated) UserID(request *http.Request) string {
	return ""
}

func (AlwaysAuthenticated) SetUserID(responseWriter http.ResponseWriter, request *http.Request, userId string) {
	// nothing to do
}

// NeverAuthenticated is an Authenticator that returns no user as being
// authenticated.
type NeverAuthenticated struct{}

// Ensure that Authenticator interface is implemented
var _ Authenticator = (*NeverAuthenticated)(nil)

func NewNeverAuthenticated() *NeverAuthenticated {
	return &NeverAuthenticated{}
}

func (NeverAuthenticated) IsAuthenticated(request *http.Request) bool {
	return false
}

func (NeverAuthenticated) UserID(request *http.Request) string {
	return ""
}

func (NeverAuthenticated) SetUserID(responseWriter http.ResponseWriter, request *http.Request, userId string) {
	// nothing to do
}
