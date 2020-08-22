package errors

// Unwrap try to find package error
func Unwrap(err error) (*Error, bool) {
	for {
		if errPkg, ok := err.(*Error); ok {
			return errPkg, true
		}
		errC, ok := err.(Causer)
		if !ok {
			return nil, false
		}
		err = errC.Cause()
	}
}
