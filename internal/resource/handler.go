package resource

import (
	"context"
	"errors"
	"log"
	"net/http"

	"nilchan-hackaton/internal/auth/token"
	"nilchan-hackaton/internal/httpapi/request"
	"nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type resourceService interface {
	Add(ctx context.Context, userID int64, request CreateResourceRequest) (*Resource, error)
}

type Handler struct {
	service   resourceService
	validator *validator.Validate
}

func NewHandler(service resourceService, validate *validator.Validate) *Handler {
	return &Handler{service: service, validator: validate}
}

func (h *Handler) HandleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, ok := request.DecodeAndValidate[CreateResourceRequest](w, r, h.validator)
		if !ok {
			return
		}
		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			renderResourceError(w, r, http.StatusUnauthorized, "invalid token")
			return
		}

		created, err := h.service.Add(r.Context(), userID, *body)
		if err != nil {
			h.handleError(w, r, err)
			return
		}

		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, response.Ok(toResponse(created)))
	}
}

func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalidURL), errors.Is(err, ErrBlockedURL):
		renderResourceError(w, r, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrDuplicate):
		renderResourceError(w, r, http.StatusConflict, ErrDuplicate.Error())
	case errors.Is(err, ErrBacklogFull):
		renderResourceError(w, r, http.StatusConflict, ErrBacklogFull.Error())
	case errors.Is(err, ErrInsufficientEPoints):
		renderResourceError(w, r, http.StatusConflict, ErrInsufficientEPoints.Error())
	case errors.Is(err, ErrContentTooShort), errors.Is(err, ErrContentTooLong):
		renderResourceError(w, r, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, ErrFirecrawlTimeout):
		renderResourceError(w, r, http.StatusGatewayTimeout, ErrFirecrawlTimeout.Error())
	case errors.Is(err, ErrFirecrawlFailed):
		renderResourceError(w, r, http.StatusBadGateway, ErrFirecrawlFailed.Error())
	case errors.Is(err, ErrStateConflict):
		renderResourceError(w, r, http.StatusConflict, ErrStateConflict.Error())
	default:
		log.Printf("resource creation failed request_id=%q error=%v", middleware.GetReqID(r.Context()), err)
		renderResourceError(w, r, http.StatusInternalServerError, "internal server error")
	}
}

func renderResourceError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, response.Error(message))
}

func toResponse(value *Resource) ResourceResponse {
	return ResourceResponse{
		ID:            value.ID,
		URL:           value.URL,
		Title:         value.Title,
		Tags:          value.Tags,
		Status:        string(value.Status),
		CreatedAt:     value.CreatedAt,
		CompletedAt:   value.CompletedAt,
		XPEarned:      value.XPEarned,
		EPointsEarned: value.EPointsEarned,
	}
}
