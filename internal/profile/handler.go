package profile

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"nilchan-hackaton/internal/auth/token"
	httpapi_response "nilchan-hackaton/internal/httpapi/response"
)

type profileService interface {
	GetProfile(
		ctx context.Context,
		userID int64,
	) (*GetProfileResponse, error)
	UpdateCosmetics(
		ctx context.Context,
		userID int64,
		req UpdateCosmeticsRequest,
	) (*GetProfileResponse, error)
}

type Handler struct {
	service profileService
}

func NewHandler(service profileService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			httpapi_response.RenderError(
				w,
				r,
				http.StatusUnauthorized,
				"invalid authentication",
			)
			return
		}

		result, err := h.service.GetProfile(r.Context(), userID)
		if err != nil {
			switch {
			case errors.Is(err, ErrProfileNotFound):
				httpapi_response.RenderError(
					w,
					r,
					http.StatusNotFound,
					"profile not found",
				)
			default:
				httpapi_response.RenderError(
					w,
					r,
					http.StatusInternalServerError,
					"failed to load profile",
				)
			}
			return
		}

		httpapi_response.RenderJSON(
			w,
			r,
			http.StatusOK,
			result,
		)
	}
}

func (h *Handler) HandleUpdateCosmetics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			httpapi_response.RenderError(
				w,
				r,
				http.StatusUnauthorized,
				"invalid authentication",
			)
			return
		}

		var request UpdateCosmeticsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpapi_response.RenderError(
				w,
				r,
				http.StatusBadRequest,
				"invalid request body",
			)
			return
		}

		response, err := h.service.UpdateCosmetics(
			r.Context(),
			userID,
			request,
		)
		if err != nil {
			switch {
			case errors.Is(err, ErrUnknownCosmetic):
				httpapi_response.RenderError(
					w,
					r,
					http.StatusUnprocessableEntity,
					"unknown cosmetic",
				)
			case errors.Is(err, ErrCosmeticTypeMismatch):
				httpapi_response.RenderError(
					w,
					r,
					http.StatusUnprocessableEntity,
					"cosmetic type does not match slot",
				)
			case errors.Is(err, ErrCosmeticNotOwned):
				httpapi_response.RenderError(
					w,
					r,
					http.StatusUnprocessableEntity,
					"cosmetic is not owned",
				)
			case errors.Is(err, ErrProfileNotFound):
				httpapi_response.RenderError(
					w,
					r,
					http.StatusNotFound,
					"profile not found",
				)
			default:
				httpapi_response.RenderError(
					w,
					r,
					http.StatusInternalServerError,
					"failed to update cosmetics",
				)
			}
			return
		}
		httpapi_response.RenderJSON(
			w,
			r,
			http.StatusOK,
			response,
		)
	}
}
