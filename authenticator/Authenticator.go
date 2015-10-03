package authenticator

import (
	"github.com/abbot/go-http-auth"
	"github.com/fxnn/gone/router"
	"net/http"
)

const (
	authenticationRealmName = "gone wiki"
)

// Authenticator provides UI and algorithms to allow an user to
// authenticate.
type Authenticator struct {
	authenticationHandler *auth.BasicAuth
}

func New() Authenticator {
	return Authenticator{auth.NewBasicAuthenticator(authenticationRealmName, provideSampleSecret)}
}

func (a *Authenticator) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if a.IsAuthenticated(request) {
		router.RedirectToViewMode(writer, request)
		return
	}

	a.authenticationHandler.RequireAuth(writer, request)
}

func (a *Authenticator) IsAuthenticated(request *http.Request) bool {
	return a.authenticationHandler.CheckAuth(request) != ""
}

func provideSampleSecret(user, realm string) string {
	if user == "test" {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}
