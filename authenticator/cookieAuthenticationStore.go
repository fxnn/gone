package authenticator

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const (
	authenticationStoreSessionName       = "goneAuthenticationStore"
	cookieAuthenticationKeyLengthInBytes = 64
)

type cookieAuthenticationStore struct {
	cookieStore sessions.Store
}

func newCookieAuthenticationStore() cookieAuthenticationStore {
	return cookieAuthenticationStore{createCookieStoreWithRandomKey()}
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

func (s *cookieAuthenticationStore) userId(request *http.Request) (string, bool) {
	var session = s.session(request)
	var userId = session.Values["userId"]
	if userId != nil {
		if strVal, ok := userId.(string); ok {
			return strVal, true
		} else {
			log.Printf("%s %s: failed to read userId from value %s", request.Method, request.URL, userId)
		}
	}
	return "", false
}

func (s *cookieAuthenticationStore) setUserId(
	writer http.ResponseWriter, request *http.Request, userId string,
) {
	var session = s.session(request)
	session.Values["userId"] = userId
	if err := s.cookieStore.Save(request, writer, session); err != nil {
		log.Printf("%s %s: failed to store userid in cookie", request.Method, request.URL)
	}
}

func (s *cookieAuthenticationStore) session(request *http.Request) *sessions.Session {
	session, err := s.cookieStore.Get(request, authenticationStoreSessionName)
	if err != nil {
		log.Printf("%s %s: failed to decode existing cookie session", request.Method, request.URL)
	}
	return session
}
