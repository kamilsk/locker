package locker

// Error defines the package errors.
type Error string

// Error returns the string representation of the error.
func (err Error) Error() string {
	return string(err)
}

// InvalidIntent is the error related to bad method call.
const InvalidIntent Error = "invalid intent"
