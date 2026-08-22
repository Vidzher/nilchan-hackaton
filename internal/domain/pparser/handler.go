package pparser

import (
	"context"
	"fmt"
	"net/http"
	"nilchan-hackaton/internal/lib/api/response"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	ParsePage(ctx context.Context, pageURL string) (string, error)
}

type Handler struct {
	service Service
}

var validatorInstance = validator.New()

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleParsePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ParsePageRequestDTO

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			fmt.Println("failed to decode request:", err)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		if err := validatorInstance.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			fmt.Println("invalid request:", validateErr.Error())

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(validateErr.Error()))

			return
		}

		markdown, err := h.service.ParsePage(r.Context(), req.URL)
		if err != nil {
			fmt.Println("failed to parse page:", err)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to parse page"))

			return
		}

		render.JSON(w, r, response.Ok(ParsePageResponseDTO{
			URL:      req.URL,
			Markdown: markdown,
		}))
	}
}
