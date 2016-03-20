package router

import (
	"net/http"
	"net/url"
)

type Mode string

// Mode-Constants give names to each mode the wiki application might be in.
// The mode names are used to identify the desired mode in a URL.
const (
	ModeView     Mode = ""
	ModeEdit          = "edit"
	ModeCreate        = "create"
	ModeLogin         = "login"
	ModeDelete        = "delete"
	ModeTemplate      = "template"
)

// To returns a URL that points to the same resource, but lets the
// wiki open it in given mode.
func To(m Mode, u *url.URL) *url.URL {
	result := *u // create a copy
	result.RawQuery = string(m)
	return &result
}

// Is returns true, iff the given request specifies to open a resource in
// the given mode.
func Is(m Mode, r *http.Request) bool {
	var ok = false

	switch m {
	case ModeView:
		ok = !Is(ModeEdit, r) && !Is(ModeCreate, r) && !Is(ModeDelete, r)
	case ModeEdit, ModeDelete, ModeCreate, ModeLogin, ModeTemplate:
		_, ok = r.Form[string(m)]
	}

	return ok
}
