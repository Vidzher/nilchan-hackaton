package leaderboard

import (
	"context"
	"log"
	"net/http"

	"nilchan-hackaton/internal/auth/token"
	"nilchan-hackaton/internal/httpapi/response"

	"github.com/go-chi/chi/v5/middleware"
)

type service interface {
	List(ctx context.Context, currentUserID int64) ([]Entry, error)
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "invalid authentication")
			return
		}

		entries, err := h.service.List(r.Context(), userID)
		if err != nil {
			log.Printf("load leaderboard failed request_id=%q error=%v", middleware.GetReqID(r.Context()), err)
			response.RenderError(w, r, http.StatusInternalServerError, "failed to load leaderboard")
			return
		}

		result := make([]EntryResponse, len(entries))
		for index, entry := range entries {
			result[index] = EntryResponse{
				Rank:      entry.Rank,
				UserID:    entry.UserID,
				Username:  entry.Username,
				XP:        entry.XP,
				Level:     entry.Level,
				AvatarID:  entry.AvatarID,
				FrameID:   entry.FrameID,
				IsCurrent: entry.IsCurrent,
			}
		}
		response.RenderJSON(w, r, http.StatusOK, response.OK(result))
	}
}
