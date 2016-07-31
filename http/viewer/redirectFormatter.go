package viewer

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/fxnn/gone/log"
	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/http/templates"
)

var winUrlRegexp = regexp.MustCompile("(?m)^URL\\=(.+)$")
var firstLineRegexp = regexp.MustCompile("(?m)^(.+)$")

type redirectFormatter struct {}

func newRedirectFormatter(l templates.Loader) redirectFormatter {
	return redirectFormatter{}
}

func (f redirectFormatter) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Warnf("%s %s: could not read from reader: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	var matches = winUrlRegexp.FindSubmatch(contents)
	if matches == nil || len(matches) < 2 {
		matches = firstLineRegexp.FindSubmatch(contents)
	}
	if matches == nil || len(matches) < 2 {
		log.Warnf("%s %s: could not find a URL in file: %v", request.Method, request.URL, matches)
		failer.ServeInternalServerError(writer, request)
		return
	}

	url, err := url.Parse(string(matches[1]))
	if err != nil {
		log.Warnf("%s %s: could not parse file contents: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	http.Redirect(writer, request, url.String(), http.StatusFound)
}
