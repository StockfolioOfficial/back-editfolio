package domain

import (
	"errors"
	"github.com/stockfolioofficial/back-editfolio/util/pointer"
)

var (

	ItemNotFound = errors.New("item not found")

	UserWrongPassword = errors.New("wrong password")

	UserSignInFailedResponse = ErrorResponse{
		ErrorCode: pointer.String("U-10"),
		Message: "unauthorized",
	}

	ServerInternalErrorResponse = ErrorResponse{
		Message: "server internal error",
	}
)

type ErrorResponse struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	Message string `json:"message"`
}