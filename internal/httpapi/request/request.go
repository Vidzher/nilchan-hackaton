package request

import (
	"fmt"
	"net/http"
	"strconv"

	"nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func PathInt64(r *http.Request, name string) (int64, error) {
	value, err := strconv.ParseInt(chi.URLParam(r, name), 10, 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return value, nil
}

func DecodeAndValidate[T any](w http.ResponseWriter, r *http.Request, validate *validator.Validate) (*T, bool) {
	var request T
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "invalid request body")
		return nil, false
	}
	if err := validate.Struct(request); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	return &request, true
}
