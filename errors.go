package rundeck

const (
	ErrCodeInvalidRequest          = "INVALID_REQUEST"
	ErrCodeDeletedRunningExecution = "RUNNING_EXECUTION"
	ErrCodeUnexpected              = "UNEXPECTED_ERROR"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
