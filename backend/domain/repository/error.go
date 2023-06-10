package repository

type NotFound struct {
	Err error
}

func (e *NotFound) Error() string {
	return e.Err.Error()
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFound)
	return ok
}
