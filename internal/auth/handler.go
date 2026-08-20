package auth

import (
	"fmt"
	"net/http"
	"nilchan-hackaton/internal/lib/api/response"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Login(email string, password string) (*AuthResult, error)
	Register(email string, password string) error
}

type Handler struct {
	service Service
}

var _validator *validator.Validate = validator.New()

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req AuthRequestDTO
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		if err := _validator.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			fmt.Println("invalid request")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(validateErr.Error()))

			return
		}

		res, err := h.service.Login(req.Email, req.Password)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error(err.Error()))

			return
		}

		render.JSON(w, r, response.Ok(
			AuthResponseDTO{
				UserId: res.user.ID,
				Email:  res.user.Email,
				Token:  res.token,
			},
		))
	}
}

func (h *Handler) HandleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req AuthRequestDTO
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		if err := _validator.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			fmt.Println("invalid request")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(validateErr.Error()))

			return
		}

		err = h.service.Register(req.Email, req.Password)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to register user"))

			return
		}
	}
}
