package deps

import "github.com/go-playground/validator/v10"

type Validator interface {
	Struct(i interface{}) error
}

type GoPlaygroundValidator struct {
	validate *validator.Validate
}

func NewGoPlaygroundValidator() *GoPlaygroundValidator {
	return &GoPlaygroundValidator{
		validate: validator.New(),
	}
}

func (v *GoPlaygroundValidator) Struct(i interface{}) error {
	return v.validate.Struct(i)
}
