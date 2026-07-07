package validator


import "github.com/go-playground/validator/v10"

type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(data any) error {
	return v.validate.Struct(data)
}

