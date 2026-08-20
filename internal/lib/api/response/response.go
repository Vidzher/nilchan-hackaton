package response

type SuccessResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func Ok(data any) SuccessResponse {
	return SuccessResponse{
		Status: StatusOk,
		Data:   data,
	}
}

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error:  msg,
	}
}
