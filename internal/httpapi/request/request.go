package request

import (
	"net/http"

	"nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func DecodeAndValidate[T any](w http.ResponseWriter, r *http.Request, validate *validator.Validate) (*T, bool) {
	var request T
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		renderError(w, r, http.StatusBadRequest, "invalid request body")
		return nil, false
	}
	if err := validate.Struct(request); err != nil {
		renderError(w, r, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	return &request, true
}

func renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, response.Error(message))
}
