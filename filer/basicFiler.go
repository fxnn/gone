package filer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The basicFiler implements basic algorithms and error handling needed when
// dealing with files.
// It basically wraps golang library functions for error handling.
type basicFiler struct {
	err error
}

// Never forget to check for errors.
// One call to this function resets the error state.
func (f *basicFiler) Err() error {
	var result = f.err
	f.err = nil
	return result
}

func (f *basicFiler) assertPathInsideWorkingDirectory(path string) {
	if f.err != nil {
		return
	}

	var normalizedPath = f.normalizePath(path)
	var wdPath = f.normalizePath(f.workingDirectory())

	if f.err == nil && !strings.HasPrefix(normalizedPath, wdPath) {
		f.err = NewPathNotFoundError(fmt.Sprintf("%s is not inside working directory", path))
	}
}

func (f *basicFiler) normalizePath(path string) string {
	if f.err != nil {
		return path
	}

	var result string

	result = f.absPath(path)

	f.assertPathExists(result)
	result = f.evalSymlinks(result)

	// HINT: Remove .. and ., remove trailing slash
	return f.cleanPath(result)
}

func (f *basicFiler) absPath(path string) (absPath string) {
	if f.err != nil {
		return path
	}
	absPath, f.err = filepath.Abs(path)
	return
}

func (f *basicFiler) assertPathExists(path string) {
	if f.err != nil {
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f.err = NewPathNotFoundError(err.Error())
	}
}

func (f *basicFiler) evalSymlinks(path string) (hardPath string) {
	if f.err != nil {
		return path
	}
	hardPath, f.err = filepath.EvalSymlinks(path)
	return
}

func (f *basicFiler) cleanPath(path string) string {
	if f.err != nil {
		return path
	}
	return filepath.Clean(path)
}

func (f *basicFiler) workingDirectory() (wd string) {
	if f.err != nil {
		return ""
	}
	wd, f.err = os.Getwd()
	return
}
