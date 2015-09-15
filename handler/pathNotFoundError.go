package handler

type pathNotFoundError string

func newPathNotFoundError(msg string) pathNotFoundError {
	return pathNotFoundError(msg)
}

func (e pathNotFoundError) Error() string {
	return string(e)
}

func isPathNotFoundError(e interface{}) bool {
	_, ok := e.(pathNotFoundError)
	return ok
}
