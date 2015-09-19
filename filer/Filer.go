package filer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Maps incoming HTTP requests to the file system.
type Filer struct{}

// Initializes a zeroe'd instance ready to use.
func NewFiler() Filer {
	return Filer{}
}

// Opens a reader for the given request. A caller must close the reader after
// using it.
func (f *Filer) OpenReader(request *http.Request) (io.ReadCloser, error) {
	var path = "." + request.URL.Path

	var err = f.assertPathInsideWorkingDirectory(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *Filer) assertPathInsideWorkingDirectory(path string) error {
	normalizedPath, err := f.normalizePath(path)
	if err != nil {
		if IsPathNotFoundError(err) {
			return err
		}
		return fmt.Errorf("checking %s inside wd: %s", path, err)
	}

	wdPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("checking %s inside wd: %s", path, err)
	}
	normalizedWdPath, err := f.normalizePath(wdPath)
	if err != nil {
		return fmt.Errorf("checking %s inside wd: %s", path, err)
	}

	if !strings.HasPrefix(normalizedPath, normalizedWdPath) {
		return NewPathNotFoundError(fmt.Sprintf("%s is not inside working directory", path))
	}

	return nil
}

func (f *Filer) normalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, fmt.Errorf("building abs path of %s: %s", path, err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return path, NewPathNotFoundError(err.Error())
	}

	hardPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return path, fmt.Errorf("removing symlinks from %s: %s", absPath, err)
	}

	// HINT: Remove .. and ., remove trailing slash
	cleanPath := filepath.Clean(hardPath)

	return cleanPath, nil
}
