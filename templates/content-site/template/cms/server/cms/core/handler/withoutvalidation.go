package handler

// HandlerWithoutExtractDataAndValidation - it is a simple Interface implementation for handlers without
// extract data from request/validation
type HandlerWithoutExtractDataAndValidation struct {
}

func (h *HandlerWithoutExtractDataAndValidation) ExtractData(_ *Context) (interface{}, error) {
	return nil, nil
}

func (h *HandlerWithoutExtractDataAndValidation) Validate(_ *Context, _ interface{}) (*Validation, error) {
	return NewValidation(), nil
}
