package errors

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.errorType.errorCode == CodeNotFound
}

func IsInternal(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.errorType.errorCode == CodeInternal
}

func IsConflict(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.errorType.errorCode == CodeConflict
}
