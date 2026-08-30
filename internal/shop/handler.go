package shop

import (
	"context"
	"errors"
	"log"
	"net/http"

	"nilchan-hackaton/internal/auth/token"
	"nilchan-hackaton/internal/cosmetics"
	"nilchan-hackaton/internal/httpapi/request"
	"nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type service interface {
	List() []cosmetics.Item
	Purchase(ctx context.Context, userID int64, itemID string) (*PurchaseResult, error)
}

type Handler struct {
	service   service
	validator *validator.Validate
}

func NewHandler(service service, validate *validator.Validate) *Handler {
	return &Handler{service: service, validator: validate}
}

func (h *Handler) HandleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := h.service.List()
		result := make([]CatalogItemResponse, len(items))
		for index, item := range items {
			result[index] = toCatalogItemResponse(item)
		}
		response.RenderJSON(w, r, http.StatusOK, response.OK(result))
	}
}

func (h *Handler) HandlePurchase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, ok := request.DecodeAndValidate[PurchaseRequest](w, r, h.validator)
		if !ok {
			return
		}

		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "invalid authentication")
			return
		}

		result, err := h.service.Purchase(r.Context(), userID, body.ItemID)
		if err != nil {
			h.handlePurchaseError(w, r, err)
			return
		}

		response.RenderJSON(w, r, http.StatusOK, response.OK(PurchaseResponse{
			Item:    toCatalogItemResponse(result.Item),
			EPoints: result.EPoints,
		}))
	}
}

func (h *Handler) handlePurchaseError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrUnknownItem):
		response.RenderError(w, r, http.StatusUnprocessableEntity, ErrUnknownItem.Error())
	case errors.Is(err, ErrItemNotPurchasable):
		response.RenderError(w, r, http.StatusUnprocessableEntity, ErrItemNotPurchasable.Error())
	case errors.Is(err, ErrAlreadyOwned):
		response.RenderError(w, r, http.StatusConflict, ErrAlreadyOwned.Error())
	case errors.Is(err, ErrInsufficientEPoints):
		response.RenderError(w, r, http.StatusConflict, ErrInsufficientEPoints.Error())
	default:
		log.Printf("purchase cosmetic failed request_id=%q error=%v", middleware.GetReqID(r.Context()), err)
		response.RenderError(w, r, http.StatusInternalServerError, "failed to purchase cosmetic")
	}
}

func toCatalogItemResponse(item cosmetics.Item) CatalogItemResponse {
	return CatalogItemResponse{
		ID:          item.ID,
		Type:        string(item.Type),
		DisplayName: item.DisplayName,
		Price:       item.Price,
		Presentation: PresentationResponse{
			Emoji:    item.Presentation.Emoji,
			AssetKey: item.Presentation.AssetKey,
			CSSClass: item.Presentation.CSSClass,
		},
	}
}
