package request

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidBody = errors.New("invalid request body")
	ErrValidation  = errors.New("invalid request")

	validate = validator.New()
)

func DecodeAndValidate[T any](r *http.Request) (*T, error) {
	var request T
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidBody, err)
	}
	if err := validate.Struct(request); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	return &request, nil
}
