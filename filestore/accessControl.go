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
		// HINT: OK, as long as the gone process can read the file
		return true
	}

	var p = a.pathFromRequest(request)
	if !a.canEnterAllParentDirectories(p) {
		return false
	}
	if !p.IsExists() {
		// HINT: Create file
		return a.canWriteDirectory(p.Dir())
	}
	return a.canWriteFile(p)
}

func (a *accessControl) HasReadAccessForRequest(request *http.Request) bool {
	if a.authenticator.IsAuthenticated(request) {
		// HINT: OK, as long as the gone process can read the file
		return true
	}

	var p = a.pathFromRequest(request)
	if !a.canEnterAllParentDirectories(p) {
		return false
	}
	return a.canReadFile(p)
}

// hasAccessForAllParentDirectories returns true iff all parent directories can
// be entered using world permissions.
func (a *accessControl) canEnterAllParentDirectories(p gopath.GoPath) bool {
	var parentDir = a.contentRoot

	// NOTE: Implicitly skips the last path component, which isn't a parent directory
	for _, component := range a.pathComponentsTo(p) {
		if !a.canEnterDirectory(parentDir) {
			return false
		}
		parentDir = parentDir.JoinPath(component)
	}

	return true
}

func (a *accessControl) canEnterDirectory(p gopath.GoPath) bool {
	if p.HasErr() || !p.IsDirectory() {
		return false
	}
	return a.hasWorldExecutePermission(p.FileMode())
}

func (a *accessControl) canWriteDirectory(p gopath.GoPath) bool {
	if p.HasErr() || !p.IsDirectory() {
		return false
	}
	return a.hasWorldWritePermission(p.FileMode())
}

func (a *accessControl) canReadFile(p gopath.GoPath) bool {
	if p.HasErr() || !p.IsRegular() {
		return false
	}
	return a.hasWorldReadPermission(p.FileMode())
}

func (a *accessControl) canWriteFile(p gopath.GoPath) bool {
	if p.HasErr() || !p.IsRegular() {
		return false
	}
	return a.hasWorldWritePermission(p.FileMode())
}

func (a *accessControl) hasWorldExecutePermission(mode os.FileMode) bool {
	return mode&0001 != 0
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
