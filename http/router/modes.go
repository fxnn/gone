package router

import (
	"net/http"
	"net/url"
)

// Mode-Constants give names to each mode the wiki application might be in.
// The mode names are used to identify the desired mode in a URL.
const (
	ModeEdit   = "edit"
	ModeCreate = "create"
	ModeLogin  = "login"
	ModeDelete = "delete"
)

// ToModeView returns a URL that points to the same resource, but lets the
// wiki open it in view mode.
func ToModeView(url *url.URL) *url.URL {
	return copyWithMode(url, "")
}

// ToModeCreate returns a URL that points to the same resource, but lets the
// wiki open it in create mode.
func ToModeCreate(url *url.URL) *url.URL {
	return copyWithMode(url, ModeCreate)
}

// ToModeEdit returns a URL that points to the same resource, but lets the
// wiki open it in edit mode.
func ToModeEdit(url *url.URL) *url.URL {
	return copyWithMode(url, ModeEdit)
}

// ToModeDelete returns a URL that points to the same resource, but lets the
// wiki open it in delete mode.
func ToModeDelete(url *url.URL) *url.URL {
	return copyWithMode(url, ModeDelete)
}

// ToModeLogin returns a URL that points to the same resource, but lets the
// wiki open it in login mode.
func ToModeLogin(url *url.URL) *url.URL {
	return copyWithMode(url, ModeLogin)
}

// Creates a copy of the given URL and changes the mode to the given value.
func copyWithMode(url *url.URL, mode string) *url.URL {
	result := *url // create a copy
	result.RawQuery = mode
	return &result
}

func IsModeView(request *http.Request) bool {
	return !IsModeEdit(request) && !IsModeCreate(request) && !IsModeDelete(request)
}

func IsModeEdit(request *http.Request) bool {
	_, ok := request.Form[ModeEdit]
	return ok
}

func IsModeDelete(request *http.Request) bool {
	_, ok := request.Form[ModeDelete]
	return ok
}

func IsModeCreate(request *http.Request) bool {
	_, ok := request.Form[ModeCreate]
	return ok
}

func IsModeLogin(request *http.Request) bool {
	_, ok := request.Form[ModeLogin]
	return ok
}
