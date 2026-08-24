package validation

import "github.com/go-playground/validator/v10"

func New() (*validator.Validate, error) {
	validate := validator.New()
	if err := validate.RegisterValidation("bcrypt_max_bytes", func(fl validator.FieldLevel) bool {
		return len([]byte(fl.Field().String())) <= 72
	}); err != nil {
		return nil, err
	}

	return validate, nil
}
