package locker

// Error defines the module errors.
type Error string

// Error returns the string representation of the error.
func (err Error) Error() string {
	return string(err)
}

// CriticalIssue is the error related to a bad module usage.
const CriticalIssue Error = "critical issue"

// Interrupted is the error related to a timeout.
const Interrupted Error = "operation interrupted"

// InvalidIntent is the error related to a bad method call.
const InvalidIntent Error = "invalid intent"
