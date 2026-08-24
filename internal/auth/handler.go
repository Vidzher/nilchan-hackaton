package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"nilchan-hackaton/internal/httpapi/request"
	"nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type service interface {
	Login(ctx context.Context, email, password string) (*Result, error)
	Register(ctx context.Context, email, username, password string) (*Result, error)
}

type Handler struct {
	service   service
	validator *validator.Validate
}

func NewHandler(service service, validate *validator.Validate) *Handler {
	return &Handler{service: service, validator: validate}
}

func (h *Handler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := request.DecodeAndValidate[LoginRequestDTO](w, r, h.validator)
		if !ok {
			return
		}

		result, err := h.service.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, ErrInvalidCredentials) {
				renderError(w, r, http.StatusUnauthorized, "invalid credentials")
				return
			}

			logAuthError(r, "login failed", err)
			renderError(w, r, http.StatusInternalServerError, "internal server error")
			return
		}

		renderAuthResult(w, r, http.StatusOK, result)
	}
}

func (h *Handler) HandleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := request.DecodeAndValidate[RegisterRequestDTO](w, r, h.validator)
		if !ok {
			return
		}

		result, err := h.service.Register(r.Context(), req.Email, req.Username, req.Password)
		if err != nil {
			if errors.Is(err, ErrUserAlreadyExists) {
				renderError(w, r, http.StatusConflict, ErrUserAlreadyExists.Error())
				return
			}

			logAuthError(r, "registration failed", err)
			renderError(w, r, http.StatusInternalServerError, "internal server error")
			return
		}

		renderAuthResult(w, r, http.StatusCreated, result)
	}
}

func renderAuthResult(w http.ResponseWriter, r *http.Request, status int, result *Result) {
	render.Status(r, status)
	render.JSON(w, r, response.Ok(AuthResponseDTO{
		UserID:   result.User.ID,
		Email:    result.User.Email,
		Username: result.User.Username,
		Token:    result.Token,
	}))
}

func renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, response.Error(message))
}

func logAuthError(r *http.Request, message string, err error) {
	log.Printf("%s request_id=%q error=%v", message, middleware.GetReqID(r.Context()), err)
}
