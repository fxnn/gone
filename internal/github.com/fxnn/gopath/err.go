package gopath

// HasErr returns true iff the GoPath is errorneous.
// It returns false otherwise, including the empty GoPath.
func (g GoPath) HasErr() bool {
	return g.err != nil
}

// Err returns the error represented by the GoPath.
// When the GoPath is not errorneous, it returns nil.
func (g GoPath) Err() error {
	return g.err
}
