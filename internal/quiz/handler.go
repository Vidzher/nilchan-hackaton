package quiz

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

type service interface {
	Get(ctx context.Context, userID, resourceID int64) (*Quiz, error)
	Complete(ctx context.Context, userID, resourceID int64, answers []Answer) (*CompletionResult, error)
}

type Handler struct {
	service   service
	validator *validator.Validate
}

func NewHandler(service service, validate *validator.Validate) *Handler {
	return &Handler{service: service, validator: validate}
}

func (h *Handler) HandleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resourceID, err := request.PathInt64(r, "resourceID")
		if err != nil {
			response.RenderError(w, r, http.StatusBadRequest, "invalid resource ID")
			return
		}
		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "invalid token")
			return
		}

		found, err := h.service.Get(r.Context(), userID, resourceID)
		if err != nil {
			h.handleError(w, r, err)
			return
		}
		response.RenderJSON(w, r, http.StatusOK, response.OK(toGetQuizResponse(found)))
	}
}

func (h *Handler) HandleComplete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resourceID, err := request.PathInt64(r, "resourceID")
		if err != nil {
			response.RenderError(w, r, http.StatusBadRequest, "invalid resource ID")
			return
		}
		body, ok := request.DecodeAndValidate[CompleteQuizRequest](w, r, h.validator)
		if !ok {
			return
		}
		userID, err := token.UserIDFromContext(r.Context())
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "invalid token")
			return
		}

		answers := make([]Answer, len(body.Answers))
		for index, answer := range body.Answers {
			answers[index] = Answer{
				QuestionIndex: answer.QuestionIndex,
				SelectedIndex: answer.SelectedIndex,
			}
		}
		result, err := h.service.Complete(r.Context(), userID, resourceID, answers)
		if err != nil {
			h.handleError(w, r, err)
			return
		}
		response.RenderJSON(w, r, http.StatusOK, response.OK(toCompleteQuizResponse(result)))
	}
}

func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.RenderError(w, r, http.StatusNotFound, ErrNotFound.Error())
	case errors.Is(err, ErrUnavailable):
		response.RenderError(w, r, http.StatusConflict, ErrUnavailable.Error())
	case errors.Is(err, ErrInvalidAnswers):
		response.RenderError(w, r, http.StatusUnprocessableEntity, ErrInvalidAnswers.Error())
	default:
		log.Printf("quiz request failed request_id=%q error=%v", middleware.GetReqID(r.Context()), err)
		response.RenderError(w, r, http.StatusInternalServerError, "internal server error")
	}
}

func toCompleteQuizResponse(value *CompletionResult) CompleteQuizResponse {
	return CompleteQuizResponse{
		Completion: CompletionDetails{
			CompletedAt:    value.CompletedAt,
			TotalQuestions: value.TotalQuestions,
			XPEarned:       value.XPEarned,
			EPointsEarned:  value.EPointsEarned,
		},
		Progress: ProgressSnapshot{
			XP:                 value.XP,
			EPoints:            value.EPoints,
			CurrentStreak:      value.CurrentStreak,
			Level:              value.Level,
			ActiveBacklogLimit: value.ActiveBacklogLimit,
		},
	}
}

func toGetQuizResponse(value *Quiz) GetQuizResponse {
	questions := make([]QuizQuestionResponse, len(value.Questions))
	for index, question := range value.Questions {
		questions[index] = QuizQuestionResponse{
			Text:              question.Text,
			Options:           question.Options,
			Explanation:       question.Explanation,
			Evidence:          question.Evidence,
			VerificationSalt:  question.VerificationSalt,
			CorrectAnswerHash: question.CorrectAnswerHash,
		}
	}
	return GetQuizResponse{
		ID:         value.ID,
		ResourceID: value.ResourceID,
		Title:      value.Title,
		Questions:  questions,
	}
}
