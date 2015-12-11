package filestore

import (
	"fmt"
	"os"

	"github.com/fxnn/gone/internal/github.com/fxnn/gopath"
	"github.com/fxnn/gone/store"
)

type errStore struct {
	err error
}

func newErrStore() *errStore {
	return &errStore{}
}

func (s *errStore) hasErr() bool {
	return s.err != nil
}

func (s *errStore) hasPathNotFoundError() bool {
	return s.hasErr() && store.IsPathNotFoundError(s.err)
}

// Never forget to check for errors.
// One call to this function resets the error state.
func (s *errStore) errAndClear() error {
	var result = s.err
	s.err = nil
	return result
}

func (s *errStore) setErr(err error) {
	s.err = s.wrapErr(err)
}

// wrapErr wraps s.err to a filestore-specific error, if possible
func (s *errStore) wrapErr(err error) error {
	if os.IsNotExist(err) {
		if pathError, ok := err.(*os.PathError); ok {
			return store.NewPathNotFoundError("path not found: " + pathError.Path)
		}
		return store.NewPathNotFoundError(fmt.Sprintf("path not found: %s", err))
	}
	return err
}

// syncedErrs couples GoPath's error handling with errStore's error handling.
// When the GoPath contained an error, it will be stored in the errStore, so
// that all following ops become no-ops.
// When however the errStore contains an error, an errorneous GoPath will be
// returned.
func (s *errStore) syncedErrs(p gopath.GoPath) gopath.GoPath {
	if s.hasErr() {
		return gopath.FromErr(s.err)
	}
	if p.HasErr() {
		s.setErr(p.Err())
	}
	return p
}
