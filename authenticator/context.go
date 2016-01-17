package authenticator

import (
	"net/http"

	"github.com/fxnn/gone/context"
)

// ContextAuthenticator saves authentication information in the rqeuest context.
type ContextAuthenticator struct {
}

func NewContextAuthenticator() *ContextAuthenticator {
	return &ContextAuthenticator{}
}

func (a *ContextAuthenticator) IsAuthenticated(request *http.Request) bool {
	return context.Load(request).IsAuthenticated()
}

func (a *ContextAuthenticator) UserID(request *http.Request) string {
	return context.Load(request).UserId
}

func (a *ContextAuthenticator) SetUserID(writer http.ResponseWriter, request *http.Request, userId string) {
	var ctx = context.Load(request)
	ctx.UserId = userId
	ctx.Save(request)
}
