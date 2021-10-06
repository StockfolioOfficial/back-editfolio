package domain

var (
	ServerInternalErrorResponse = ErrorResponse{
		Message: "server internal error",
	}
)

type ErrorResponse struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	Message string `json:"message"`
}