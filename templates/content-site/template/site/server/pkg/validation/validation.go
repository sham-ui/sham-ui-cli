package validation

type Validation struct {
	Errors []string
}

func (v *Validation) AddErrors(message ...string) *Validation {
	v.Errors = append(v.Errors, message...)
	return v
}

func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

func New() *Validation {
	return &Validation{
		Errors: nil,
	}
}
