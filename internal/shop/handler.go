package shop

import (
	"net/http"

	"nilchan-hackaton/internal/cosmetics"
	"nilchan-hackaton/internal/httpapi/response"
)

type service interface {
	List() []cosmetics.Item
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := h.service.List()
		result := make([]CatalogItemResponse, len(items))
		for index, item := range items {
			result[index] = CatalogItemResponse{
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
		response.RenderJSON(w, r, http.StatusOK, response.OK(result))
	}
}
