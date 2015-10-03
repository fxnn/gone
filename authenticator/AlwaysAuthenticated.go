package authenticator

import "net/http"

type AlwaysAuthenticated struct{}

func NewAlwaysAuthenticated() *AlwaysAuthenticated {
	return &AlwaysAuthenticated{}
}

func (AlwaysAuthenticated) IsAuthenticated(request *http.Request) bool {
	return true
}
