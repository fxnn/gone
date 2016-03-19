package authenticator

import (
	"net/http"
	"time"

	"github.com/abbot/go-http-auth"
	"github.com/fxnn/gone/authenticator/bruteblocker"
	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/http/router"
	"github.com/fxnn/gone/log"
	"github.com/fxnn/gopath"
)

const (
	authenticationRealmName = "gone wiki"
)

// HttpBasicAuthenticator is an HttpAuthenticator that uses HTTP Basic Auth for
// initial authentication and stores the result in a session cookie for further
// requests.
type HttpBasicAuthenticator struct {
	requestAuth         Authenticator // requestAuth stores authentication information during this request.
	sessionAuth         Authenticator // sessionAuth stores authentication information during the user session.
	loginRequiresHeader string
	basicAuth           *auth.BasicAuth
	bruteBlocker        *bruteblocker.BruteBlocker
}

// NewHttpBasicAuthenticator creates a new instance.
//
// requestAuth will be provided with the auth information for each request.
// htpasswdFile is used as source of usernames and passwords.
// loginRequiresHeader is the name of an HTTP header required for each login
// attempt.
// This may be used to only allow login over secured connections.
// bruteBlocker is a configured BruteBlocker instance.
func NewHttpBasicAuthenticator(
	requestAuth Authenticator,
	htpasswdFile gopath.GoPath,
	loginRequiresHeader string,
	bruteBlocker *bruteblocker.BruteBlocker,
) *HttpBasicAuthenticator {
	return &HttpBasicAuthenticator{
		NewContextAuthenticator(),
		NewCookieAuthenticator(),
		loginRequiresHeader,
		createBasicAuth(authenticationRealmName, htpasswdFile),
		bruteBlocker}
}

func createBasicAuth(realmName string, htpasswdFile gopath.GoPath) *auth.BasicAuth {
	var secretProvider = noSecrets
	if !htpasswdFile.HasErr() && !htpasswdFile.IsEmpty() {
		secretProvider = auth.HtpasswdFileProvider(htpasswdFile.Path())
	}
	return auth.NewBasicAuthenticator(realmName, secretProvider)
}

func (a *HttpBasicAuthenticator) MiddlewareHandler(delegate http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// copy from session cookie
		if userID := a.sessionAuth.UserID(request); userID != "" {
			a.requestAuth.SetUserID(writer, request, userID)
		}

		delegate.ServeHTTP(writer, request)
	})
}

func (a *HttpBasicAuthenticator) LoginHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if a.loginRequiresHeader != "" && request.Header.Get(a.loginRequiresHeader) == "" {
			log.Printf("%s %s: deny login because of missing connection header '%s'",
				request.Method, request.URL, a.loginRequiresHeader)
			failer.ServeBadRequest(writer, request)
			return
		}

		var user = a.userAttemptingAuth(request)
		if user != "" {
			a.authenticate(writer, request)

			// NOTE: Delay request even if authentication was successful, so that the
			// attacker needs our response
			time.Sleep(a.bruteBlocker.Delay(user, request.RemoteAddr, a.requestAuth.IsAuthenticated(request)))

			if a.requestAuth.IsAuthenticated(request) && a.requestAuth.UserID(request) == user {
				log.Printf("%s %s: authenticated as %s", request.Method, request.URL, a.requestAuth.UserID(request))
				a.sessionAuth.SetUserID(writer, request, a.requestAuth.UserID(request))
				router.RedirectToViewMode(writer, request)
				return
			}
		}

		a.basicAuth.RequireAuth(writer, request)
	})
}

func (a *HttpBasicAuthenticator) userAttemptingAuth(request *http.Request) string {
	if user, _, ok := request.BasicAuth(); ok {
		return user
	}
	return ""
}

func (a *HttpBasicAuthenticator) authenticate(writer http.ResponseWriter, request *http.Request) {
	a.requestAuth.SetUserID(writer, request, a.basicAuth.CheckAuth(request))
}
