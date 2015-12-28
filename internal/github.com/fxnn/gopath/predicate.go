package gopath

import "os"

// IsExists returns true iff we have no error and Stat() succeeds.
func (g GoPath) IsExists() bool {
	return !g.AssertExists().HasErr()
}

// IsDirectory returns true iff we have no error and the path points to a directory.
func (g GoPath) IsDirectory() bool {
	if info := g.FileInfo(); info != nil {
		return info.IsDir()
	} else {
		return false
	}
}

func (g GoPath) HasFileMode(mode os.FileMode) bool {
	if g.Stat(); !g.HasErr() {
		return g.FileMode()&mode != 0
	}
	return false
}

func (g GoPath) IsRegular() bool {
	if g.Stat(); !g.HasErr() {
		// NOTE: Extra precaution has to be taken, as FileMode() defaults to 0,
		// which causes IsRegular() to return true ...
		return g.FileMode().IsRegular()
	}
	return false
}

func (g GoPath) IsSymlink() bool {
	return g.HasFileMode(os.ModeSymlink)
}

func (g GoPath) IsTemporary() bool {
	return g.HasFileMode(os.ModeTemporary)
}
