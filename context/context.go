package context

import (
	"github.com/gorilla/context"
	"net/http"
)

const (
	userIdKey = iota
)

type Context struct {
	// UserId is the unique id of the user, when he authenticated, or the empty
	// string otherwise.
	UserId string
}

func Load(request *http.Request) Context {
	var result = Context{}
	if val, ok := context.GetOk(request, userIdKey); ok {
		result.UserId = val.(string)
	}
	return result
}

func (c Context) Save(request *http.Request) {
	context.Clear(request)
	if c.UserId != "" {
		context.Set(request, userIdKey, c.UserId)
	}
}

func (c Context) IsAuthenticated() bool {
	return c.UserId != ""
}
