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

func (a *accessControl) assertHasWriteAccessForRequest(request *http.Request) {
	if a.err != nil {
		return
	}
	if !a.HasWriteAccessForRequest(request) {
		a.setErr(NewAccessDeniedError(fmt.Sprintf("Access denied on %s", request.URL)))
	}
}

func (a *accessControl) assertHasReadAccessForRequest(request *http.Request) {
	if a.err != nil {
		return
	}
	if !a.HasReadAccessForRequest(request) {
		a.setErr(NewAccessDeniedError(fmt.Sprintf("Access denied on %s", request.URL)))
	}
}

func (a *accessControl) HasWriteAccessForRequest(request *http.Request) bool {
	if a.authenticator.IsAuthenticated(request) {
		return true
	}
	return a.hasWriteAccessToPath(a.pathFromRequest(request))
}

func (a *accessControl) HasReadAccessForRequest(request *http.Request) bool {
	if a.authenticator.IsAuthenticated(request) {
		return true
	}
	return a.hasReadAccessToPath(a.pathFromRequest(request))
}

func (a *accessControl) hasWriteAccessToPath(p string) bool {
	if a.err != nil {
		return false
	}
	return a.hasWriteAccessForFileMode(a.relevantFileModeForPath(p))
}

func (a *accessControl) hasReadAccessToPath(p string) bool {
	if a.err != nil {
		return false
	}
	return a.hasReadAccessForFileMode(a.relevantFileModeForPath(p))
}

func (a *accessControl) hasWriteAccessForFileMode(mode os.FileMode) bool {
	// 0002 is the write permission for others
	return mode&0002 != 0
}

func (a *accessControl) hasReadAccessForFileMode(mode os.FileMode) bool {
	// 0004 is the read permission for others
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
