package filer

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/fxnn/gone/authenticator"
)

// Maps incoming HTTP requests to the file system.
type accessControl struct {
	authenticator authenticator.Authenticator
	basicFiler
}

func newAccessControl(authenticator authenticator.Authenticator) accessControl {
	return accessControl{authenticator, newBasicFiler()}
}

func (a *accessControl) SetAuthenticator(authenticator authenticator.Authenticator) {
	a.authenticator = authenticator
}

func (a *accessControl) assertHasWriteAccessForRequest(request *http.Request) {
	if a.err != nil {
		return
	}
	if !a.HasWriteAccessForRequest(request) {
		var msg = fmt.Sprintf("Write access denied on %s", request.URL)
		if a.err != nil {
			msg = fmt.Sprintf("%s: %s", msg, a.err)
		}
		a.setErr(NewAccessDeniedError(msg))
	}
}

func (a *accessControl) assertHasReadAccessForRequest(request *http.Request) {
	if a.err != nil {
		return
	}
	if !a.HasReadAccessForRequest(request) {
		var msg = fmt.Sprintf("Read access denied on %s", request.URL)
		if a.err != nil {
			msg = fmt.Sprintf("%s: %s", msg, a.err)
		}
		a.setErr(NewAccessDeniedError(msg))
	}
}

func (a *accessControl) HasWriteAccessForRequest(request *http.Request) bool {
	if a.authenticator.IsAuthenticated(request) {
		return true
	}
	var mode = a.relevantFileModeForPath(a.pathFromRequest(request))
	return a.err == nil && a.hasWorldWritePermission(mode)
}

func (a *accessControl) HasReadAccessForRequest(request *http.Request) bool {
	if a.authenticator.IsAuthenticated(request) {
		return true
	}
	var mode = a.relevantFileModeForPath(a.pathFromRequest(request))
	return a.err == nil && a.hasWorldReadPermission(mode)
}

func (a *accessControl) hasWorldWritePermission(mode os.FileMode) bool {
	return mode&0002 != 0
}

func (a *accessControl) hasWorldReadPermission(mode os.FileMode) bool {
	return mode&0004 != 0
}

// getRelevantFileModeForPath returns the FileMode for the given file or, when
// the file does not exist, its containing directory.
func (a *accessControl) relevantFileModeForPath(p string) os.FileMode {
	if a.err != nil {
		return 0
	}
	info, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		// HINT: Inspect permissions of containing directory
		info, err = os.Stat(path.Dir(p))
	}
	if err != nil {
		a.setErr(err)
		return 0
	}
	return info.Mode()
}
