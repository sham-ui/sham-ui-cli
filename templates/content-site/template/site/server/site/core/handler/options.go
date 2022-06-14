package handler

type handlerOptions struct {
	serializeResultToJSON bool
}

type Option interface {
	apply(options *handlerOptions)
}

// funcHandlerOption wraps a function that modifies handlerOptions into an
// implementation of the Option interface.
type funcHandlerOption struct {
	f func(options *handlerOptions)
}

func (fho *funcHandlerOption) apply(options *handlerOptions) {
	fho.f(options)
}

func newFuncHandlerOption(f func(options *handlerOptions)) *funcHandlerOption {
	return &funcHandlerOption{
		f: f,
	}
}

func defaultHandlerOptions() handlerOptions {
	return handlerOptions{
		serializeResultToJSON: true,
	}
}

// WithoutSerializeResultToJSON disable serialization Process response to JSON
func WithoutSerializeResultToJSON() Option {
	return newFuncHandlerOption(func(options *handlerOptions) {
		options.serializeResultToJSON = false
	})
}
