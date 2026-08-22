package auth

import (
	"errors"
	"log"
	"net/http"
	"nilchan-hackaton/internal/lib/api/request"
	"nilchan-hackaton/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Service interface {
	Login(email, password string) (*AuthResult, error)
	Register(email, username, password string) (*AuthResult, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := request.DecodeAndValidate[LoginRequestDTO](r)
		if err != nil {
			renderRequestError(w, r, err)
			return
		}

		result, err := h.service.Login(req.Email, req.Password)
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
		req, err := request.DecodeAndValidate[RegisterRequestDTO](r)
		if err != nil {
			renderRequestError(w, r, err)
			return
		}

		result, err := h.service.Register(req.Email, req.Username, req.Password)
		if err != nil {
			switch {
			case errors.Is(err, ErrEmailTaken):
				renderError(w, r, http.StatusConflict, ErrEmailTaken.Error())
			case errors.Is(err, ErrUsernameTaken):
				renderError(w, r, http.StatusConflict, ErrUsernameTaken.Error())
			default:
				logAuthError(r, "registration failed", err)
				renderError(w, r, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		renderAuthResult(w, r, http.StatusCreated, result)
	}
}

func renderRequestError(w http.ResponseWriter, r *http.Request, err error) {
	message := request.ErrValidation.Error()
	if errors.Is(err, request.ErrInvalidBody) {
		message = request.ErrInvalidBody.Error()
	}
	renderError(w, r, http.StatusBadRequest, message)
}

func renderAuthResult(w http.ResponseWriter, r *http.Request, status int, result *AuthResult) {
	render.Status(r, status)
	render.JSON(w, r, response.Ok(AuthResponseDTO{
		UserID:   result.user.ID,
		Email:    result.user.Email,
		Username: result.user.Username,
		Token:    result.token,
	}))
}

func renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, response.Error(message))
}

func logAuthError(r *http.Request, message string, err error) {
	log.Printf("%s request_id=%q error=%v", message, middleware.GetReqID(r.Context()), err)
}
