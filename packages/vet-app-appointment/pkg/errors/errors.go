package errors

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(error string) *ErrorResponse {
	return &ErrorResponse{Error: error}
}
