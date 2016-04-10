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
	ModeUpload        = "upload"
)

var allModes = []Mode{
	ModeView,
	ModeEdit,
	ModeCreate,
	ModeLogin,
	ModeDelete,
	ModeTemplate,
	ModeUpload,
}

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
	if m != ModeView {
		_, ok := r.Form[string(m)]
		return ok
	}

	for _, mode := range allModes {
		if mode != ModeView && Is(mode, r) {
			return false
		}
	}

	return true
}
