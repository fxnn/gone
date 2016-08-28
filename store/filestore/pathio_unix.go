// +build linux

package filestore

import (
	"github.com/fxnn/gopath"
	"golang.org/x/sys/unix"
)

func isPathWriteable(p gopath.GoPath) bool {
	return unix.Access(p.Path(), unix.W_OK) == nil
}
