package httpapi_request

import (
	"net/http"
	httpapi_response "nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func DecodeAndValidate[T any](w http.ResponseWriter, r *http.Request, validate *validator.Validate) (*T, bool) {
	var request T
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		httpapi_response.RenderError(w, r, http.StatusBadRequest, "invalid request body")
		return nil, false
	}
	if err := validate.Struct(request); err != nil {
		httpapi_response.RenderError(w, r, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	return &request, true
}
