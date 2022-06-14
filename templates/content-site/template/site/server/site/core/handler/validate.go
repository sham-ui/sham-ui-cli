package handler

type Validation struct {
	IsValid bool
	Errors  []string
}

func (v *Validation) AddError(msg string) *Validation {
	v.IsValid = false
	v.Errors = append(v.Errors, msg)
	return v
}

func NewValidation() *Validation {
	return &Validation{
		IsValid: true,
	}
}
