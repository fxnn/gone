package authenticator

import (
	"log"
	"net/http"

	"github.com/abbot/go-http-auth"
	"github.com/fxnn/gone/context"
	"github.com/fxnn/gone/router"
)

const (
	authenticationRealmName = "gone wiki"
)

type HttpBasicAuthenticator struct {
	authenticationHandler *auth.BasicAuth
}

func NewHttpBasicAuthenticator() *HttpBasicAuthenticator {
	return &HttpBasicAuthenticator{auth.NewBasicAuthenticator(authenticationRealmName, provideSampleSecret)}
}

func (a *HttpBasicAuthenticator) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.checkAuth(request)

	if a.IsAuthenticated(request) {
		log.Printf("%s %s: authenticated as %s", request.Method, request.URL, a.UserId(request))
		router.RedirectToViewMode(writer, request)
		return
	}

	a.authenticationHandler.RequireAuth(writer, request)
}

func (a *HttpBasicAuthenticator) IsAuthenticated(request *http.Request) bool {
	return context.Load(request).IsAuthenticated()
}

func (a *HttpBasicAuthenticator) UserId(request *http.Request) string {
	return context.Load(request).UserId
}

func (a *HttpBasicAuthenticator) checkAuth(request *http.Request) {
	var ctx = context.Load(request)
	ctx.UserId = a.authenticationHandler.CheckAuth(request)
	ctx.Save(request)
}

func provideSampleSecret(user, realm string) string {
	if user == "test" {
		// password is "hello"
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}
