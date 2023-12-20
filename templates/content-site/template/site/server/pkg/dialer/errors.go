package dialer

type (
	DialUnixSocketError struct {
		Err error
	}
	DialTCPError struct {
		Err error
	}
)

func (e DialUnixSocketError) Error() string {
	return "dial unix: " + e.Err.Error()
}

func (e DialUnixSocketError) Unwrap() error {
	return e.Err
}

func (e DialUnixSocketError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(DialUnixSocketError)
	return ok
}

func (e DialTCPError) Error() string {
	return "dial tcp: " + e.Err.Error()
}

func (e DialTCPError) Unwrap() error {
	return e.Err
}

func (e DialTCPError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(DialTCPError)
	return ok
}

func NewDialUnixSocketError(err error) DialUnixSocketError {
	return DialUnixSocketError{Err: err}
}

func NewDialTCPError(err error) DialTCPError {
	return DialTCPError{Err: err}
}
