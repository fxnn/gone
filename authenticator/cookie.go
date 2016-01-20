package authenticator

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/fxnn/gone/log"
)

const (
	authenticationStoreSessionName       = "goneAuthenticationStore"
	cookieAuthenticationKeyLengthInBytes = 64
	cookieMaxAge                         = 1 * time.Hour
	userIDKey                            = "userID"
)

// CookieAuthenticator stores authentication information in a cookie in the
// user agent.
type CookieAuthenticator struct {
	cookieStore sessions.Store
}

func NewCookieAuthenticator() *CookieAuthenticator {
	var cookieStore = createCookieStoreWithRandomKey()
	cookieStore.MaxAge(int(cookieMaxAge / time.Second))
	return &CookieAuthenticator{cookieStore}
}

func createCookieStoreWithRandomKey() *sessions.CookieStore {
	var authenticationKey = securecookie.GenerateRandomKey(cookieAuthenticationKeyLengthInBytes)
	if authenticationKey == nil {
		log.Fatalf(
			"failed to generate random cookie authentication key of %d bytes",
			cookieAuthenticationKeyLengthInBytes)
	}
	return sessions.NewCookieStore(authenticationKey)
}

func (s *CookieAuthenticator) IsAuthenticated(request *http.Request) bool {
	return s.UserID(request) != ""
}

func (s *CookieAuthenticator) UserID(request *http.Request) string {
	var session = s.session(request)
	var userId = session.Values[userIDKey]
	if userId != nil {
		if strVal, ok := userId.(string); ok {
			return strVal
		} else {
			log.Printf("%s %s: failed to read userId from value %s", request.Method, request.URL, userId)
		}
	}
	return ""
}

func (s *CookieAuthenticator) SetUserID(writer http.ResponseWriter, request *http.Request, userId string) {
	if userId == "" {
		s.removeCookie(writer, request)
	} else {
		var session = s.session(request)
		session.Values[userIDKey] = userId
		if err := s.cookieStore.Save(request, writer, session); err != nil {
			log.Printf("%s %s: failed to store userid in cookie", request.Method, request.URL)
		}
	}
}

func (s *CookieAuthenticator) removeCookie(writer http.ResponseWriter, request *http.Request) {
	var session = s.session(request)
	session.Options.MaxAge = -1
	if err := s.cookieStore.Save(request, writer, session); err != nil {
		log.Printf("%s %s: failed to clear cookie", request.Method, request.URL)
	}
}

func (s *CookieAuthenticator) session(request *http.Request) *sessions.Session {
	session, err := s.cookieStore.Get(request, authenticationStoreSessionName)
	if err != nil {
		log.Printf("%s %s: failed to decode existing cookie session", request.Method, request.URL)
	}
	return session
}
