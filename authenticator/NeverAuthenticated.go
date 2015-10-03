package authenticator

import "net/http"

type NeverAuthenticated struct{}

func NewNeverAuthenticated() *NeverAuthenticated {
	return &NeverAuthenticated{}
}

func (NeverAuthenticated) IsAuthenticated(request *http.Request) bool {
	return false
}
