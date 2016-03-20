package router

import (
	"net/http"
	"net/url"
)

func RedirectToViewMode(writer http.ResponseWriter, request *http.Request) {
	Redirect(writer, request, To(ModeView, request.URL))
}

func RedirectToEditMode(writer http.ResponseWriter, request *http.Request) {
	Redirect(writer, request, To(ModeEdit, request.URL))
}

func Redirect(writer http.ResponseWriter, request *http.Request, location *url.URL) {
	http.Redirect(writer, request, location.String(), http.StatusFound)
}
