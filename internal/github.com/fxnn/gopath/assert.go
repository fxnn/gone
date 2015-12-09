package gopath

// AssertExists returns an errorneous path, iff a Stat() call fails, and sets
// the error value. Otherwise, it returns the GoPath itself.
//
// Note that it seems wrong to only return an errorneous path when Stat() fails
// with an error that yields os.IsNotExist(err) == true, because a different
// error does not necessarily mean that the file does exist.
func (g GoPath) AssertExists() GoPath {
	return g.Stat()
}
