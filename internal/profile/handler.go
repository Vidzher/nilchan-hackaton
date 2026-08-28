package profile

import (
	"context"
	"errors"
	"log"
	"net/http"

	"nilchan-hackaton/internal/auth/token"
	"nilchan-hackaton/internal/httpapi/request"
	"nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type profileService interface {
	GetProfile(ctx context.Context, userID int64) (*GetProfileResponse, error)
	UpdateCosmetics(ctx context.Context, userID int64, req UpdateCosmeticsRequest) (*GetProfileResponse, error)
}

type Handler struct {
	service   profileService
	validator *validator.Validate
}

func NewHandler(service profileService, validate *validator.Validate) *Handler {
	return &Handler{service: service, validator: validate}
}

func (h *Handler) HandleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "invalid authentication")
			return
		}

		result, err := h.service.GetProfile(r.Context(), userID)
		if err != nil {
			if errors.Is(err, ErrProfileNotFound) {
				response.RenderError(w, r, http.StatusNotFound, "profile not found")
				return
			}
			logProfileError(r, "load profile failed", err)
			response.RenderError(w, r, http.StatusInternalServerError, "failed to load profile")
			return
		}

		response.RenderJSON(w, r, http.StatusOK, response.Ok(result))
	}
}

func (h *Handler) HandleUpdateCosmetics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, ok := request.DecodeAndValidate[UpdateCosmeticsRequest](w, r, h.validator)
		if !ok {
			return
		}

		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "invalid authentication")
			return
		}

		result, err := h.service.UpdateCosmetics(r.Context(), userID, *body)
		if err != nil {
			switch {
			case errors.Is(err, ErrUnknownCosmetic):
				response.RenderError(w, r, http.StatusUnprocessableEntity, "unknown cosmetic")
			case errors.Is(err, ErrCosmeticTypeMismatch):
				response.RenderError(w, r, http.StatusUnprocessableEntity, "cosmetic type does not match slot")
			case errors.Is(err, ErrCosmeticNotOwned):
				response.RenderError(w, r, http.StatusUnprocessableEntity, "cosmetic is not owned")
			case errors.Is(err, ErrProfileNotFound):
				response.RenderError(w, r, http.StatusNotFound, "profile not found")
			default:
				logProfileError(r, "update profile cosmetics failed", err)
				response.RenderError(w, r, http.StatusInternalServerError, "failed to update cosmetics")
			}
			return
		}

		response.RenderJSON(w, r, http.StatusOK, response.Ok(result))
	}
}

func logProfileError(r *http.Request, message string, err error) {
	log.Printf("%s request_id=%q error=%v", message, middleware.GetReqID(r.Context()), err)
}
