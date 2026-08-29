package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type SuccessResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(data any) SuccessResponse {
	return SuccessResponse{
		Status: StatusOK,
		Data:   data,
	}
}

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error:  msg,
	}
}

func RenderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, Error(message))
}

func RenderJSON(w http.ResponseWriter, r *http.Request, status int, body any) {
	render.Status(r, status)
	render.JSON(w, r, body)
}
