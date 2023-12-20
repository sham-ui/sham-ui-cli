package cms

type GRPCError struct {
	Method string
	Err    error
}

func (e GRPCError) Error() string {
	return "grpc error: " + e.Method + ": " + e.Err.Error()
}

func (e GRPCError) Unwrap() error {
	return e.Err
}

func (e GRPCError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(GRPCError)
	return ok
}

func NewGRPCError(method string, err error) GRPCError {
	return GRPCError{
		Method: method,
		Err:    err,
	}
}
