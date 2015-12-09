package store

type AccessDeniedError string

func NewAccessDeniedError(msg string) AccessDeniedError {
	return AccessDeniedError(msg)
}

func (e AccessDeniedError) Error() string {
	return string(e)
}

func IsAccessDeniedError(e interface{}) bool {
	_, ok := e.(AccessDeniedError)
	return ok
}
