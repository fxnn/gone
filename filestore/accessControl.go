package filestore

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fxnn/gone/authenticator"
	"github.com/fxnn/gone/internal/github.com/fxnn/gopath"
	"github.com/fxnn/gone/store"
)

// accessControl implements permission checking for incoming requests
// based on the file system's permissions.
type accessControl struct {
	authenticator authenticator.Authenticator
	*pathIO
	*errStore
}

func newAccessControl(a authenticator.Authenticator, p *pathIO, s *errStore) *accessControl {
	return &accessControl{a, p, s}
}

func (a *accessControl) assertHasWriteAccessForRequest(request *http.Request) {
	if a.hasErr() {
		return
	}
	if !a.HasWriteAccessForRequest(request) {
		var msg = fmt.Sprintf("Write access denied on %s", request.URL)
		if a.hasErr() {
			msg = fmt.Sprintf("%s: %s", msg, a.err)
		}
		a.setErr(store.NewAccessDeniedError(msg))
	}
}

func (a *accessControl) assertHasReadAccessForRequest(request *http.Request) {
	if a.hasErr() {
		return
	}
	if !a.HasReadAccessForRequest(request) {
		var msg = fmt.Sprintf("Read access denied on %s", request.URL)
		if a.hasErr() {
			msg = fmt.Sprintf("%s: %s", msg, a.err)
		}
		a.setErr(store.NewAccessDeniedError(msg))
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
func (a *accessControl) relevantFileModeForPath(p gopath.GoPath) os.FileMode {
	if a.hasErr() || p.HasErr() {
		return 0
	}
	var pStat = p.Stat()
	if !pStat.IsExists() {
		// HINT: Inspect permissions of containing directory
		pStat = p.Dir().Stat()
	}
	a.setErr(pStat.Err())
	return pStat.FileMode()
}
