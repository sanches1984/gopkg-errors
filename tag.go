package errors

// ErrorTag ...
type ErrorTag struct {
	_ interface{} // hack empty struct
}

// Predefined tags
var (
	ApplicationError = NewTag()
	LogicError       = NewTag()
)

// NewTag ...
func NewTag() *ErrorTag {
	return &ErrorTag{}
}

// IsTagged ...
func (et *ErrorTag) IsTagged(err error) bool {
	if err == nil {
		return false
	}
	if appErr, ok := err.(*Error); ok {
		return appErr.IsTagged(et)
	}
	return false
}
