package store

type PathNotFoundError string

func NewPathNotFoundError(msg string) PathNotFoundError {
	return PathNotFoundError(msg)
}

func (e PathNotFoundError) Error() string {
	return string(e)
}

func IsPathNotFoundError(e interface{}) bool {
	_, ok := e.(PathNotFoundError)
	return ok
}
