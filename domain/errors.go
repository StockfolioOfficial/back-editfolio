package domain

import (
	"errors"

	"github.com/stockfolioofficial/back-editfolio/util/pointer"
)

var (
	ItemNotFound = errors.New("item not found")

	UserWrongPassword = errors.New("wrong password")

	UserNotAdmin = errors.New("not admin")

	InvalidateTokenResponse = ErrorResponse{
		ErrorCode: pointer.String("A-1"),
		Message:   "unauthorized",
	}

	UserSignInFailedResponse = ErrorResponse{
		ErrorCode: pointer.String("U-1"),
		Message:   "unauthorized",
	}
	
	UserWrongPasswordToUpdatePassword = ErrorResponse{
		ErrorCode: pointer.String("U-2"),
		Message:   UserWrongPassword.Error(),
	}

	ServerInternalErrorResponse = ErrorResponse{
		Message: "server internal error",
	}
)

type ErrorResponse struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	Message   string  `json:"message"`
}
