package router

import (
	"net/http"
	"net/url"
)

const (
	ModeEdit   = "edit"
	ModeCreate = "create"
)

// Returns a URL that points to the same resource, but lets the wiki open it in
// view mode.
func ToModeView(url *url.URL) *url.URL {
	return copyWithMode(url, "")
}

// Returns a URL that points to the same resource, but lets the wiki open it in
// create mode.
func ToModeCreate(url *url.URL) *url.URL {
	return copyWithMode(url, ModeCreate)
}

// Returns a URL that points to the same resource, but lets the wiki open it in
// edit mode.
func ToModeEdit(url *url.URL) *url.URL {
	return copyWithMode(url, ModeEdit)
}

// Creates a copy of the given URL and changes the mode to the given value.
func copyWithMode(url *url.URL, mode string) *url.URL {
	result := *url // create a copy
	result.RawQuery = mode
	return &result
}

func IsModeView(request *http.Request) bool {
	return !IsModeEdit(request) && !IsModeCreate(request)
}

func IsModeEdit(request *http.Request) bool {
	_, ok := request.Form[ModeEdit]
	return ok
}

func IsModeCreate(request *http.Request) bool {
	_, ok := request.Form[ModeCreate]
	return ok
}
