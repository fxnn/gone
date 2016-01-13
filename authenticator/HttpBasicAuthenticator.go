package authenticator

import (
	"encoding/base64"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/abbot/go-http-auth"
	"github.com/fxnn/gone/authenticator/bruteblocker"
	"github.com/fxnn/gone/context"
	"github.com/fxnn/gone/http/router"
	"github.com/fxnn/gopath"
)

const (
	authenticationRealmName = "gone wiki"
)

var (
	authTokenRegexp *regexp.Regexp
	authUserRegexp  *regexp.Regexp
)

func init() {
	var err error
	if authTokenRegexp, err = regexp.Compile("^[^ ]* (.*)$"); err != nil {
		panic(err)
	}
	if authUserRegexp, err = regexp.Compile("^([^:]*):"); err != nil {
		panic(err)
	}
}

func authAttemptUser(request *http.Request) string {
	var tokens = authTokenRegexp.FindStringSubmatch(request.Header.Get("Authorization"))
	if len(tokens) < 2 {
		return ""
	}
	var token = tokens[1]

	userAndPass, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return ""
	}

	var users = authUserRegexp.FindSubmatch(userAndPass)
	if len(users) < 2 {
		return ""
	}
	return string(users[1])
}

type HttpBasicAuthenticator struct {
	authenticationHandler *auth.BasicAuth
	authenticationStore   cookieAuthenticationStore
	bruteBlocker          *bruteblocker.BruteBlocker
}

func NewHttpBasicAuthenticator(
	htpasswdFile gopath.GoPath,
	delayMax time.Duration,
	userDelayStep time.Duration,
	addrDelayStep time.Duration,
	globalDelayStep time.Duration,
) *HttpBasicAuthenticator {
	var secretProvider = noSecrets
	if !htpasswdFile.HasErr() && !htpasswdFile.IsEmpty() {
		secretProvider = auth.HtpasswdFileProvider(htpasswdFile.Path())
	}
	var authenticationHandler = auth.NewBasicAuthenticator(authenticationRealmName, secretProvider)
	var authenticationStore = newCookieAuthenticationStore()
	var bruteBlocker = bruteblocker.New(delayMax, userDelayStep, addrDelayStep, globalDelayStep)
	return &HttpBasicAuthenticator{authenticationHandler, authenticationStore, bruteBlocker}
}

func (a *HttpBasicAuthenticator) AuthHandler(delegate http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if userId, ok := a.authenticationStore.userId(request); ok {
			a.setUserId(request, userId)
		} else {
			a.authenticationStore.clearUserId(writer, request)
		}

		delegate.ServeHTTP(writer, request)
	})
}

func (a *HttpBasicAuthenticator) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var user = authAttemptUser(request)
	if user != "" {
		a.checkAuth(request)

		// NOTE: Delay request even if authentication was successful, so that the
		// attacker needs our response
		time.Sleep(a.bruteBlocker.Delay(user, request.RemoteAddr, a.IsAuthenticated(request)))

		if a.IsAuthenticated(request) {
			log.Printf("%s %s: authenticated as %s", request.Method, request.URL, a.UserId(request))
			a.authenticationStore.setUserId(writer, request, a.UserId(request))
			router.RedirectToViewMode(writer, request)
			return
		}
	}

	a.authenticationHandler.RequireAuth(writer, request)
}

func (a *HttpBasicAuthenticator) checkAuth(request *http.Request) {
	a.setUserId(request, a.authenticationHandler.CheckAuth(request))
}

func (a *HttpBasicAuthenticator) IsAuthenticated(request *http.Request) bool {
	return context.Load(request).IsAuthenticated()
}

func (a *HttpBasicAuthenticator) UserId(request *http.Request) string {
	return context.Load(request).UserId
}

func (a *HttpBasicAuthenticator) setUserId(request *http.Request, userId string) {
	var ctx = context.Load(request)
	ctx.UserId = userId
	ctx.Save(request)
}

// noSecrets is a auth.SecretProvider that always fails authentication.
func noSecrets(user, realm string) string {
	// NOTE "Returning an empty string means failing the authentication."
	// (from godoc.org/github.com/abbot/go-http-auth)
	return ""
}
