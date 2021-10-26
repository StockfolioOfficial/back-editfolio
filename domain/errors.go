package domain

import (
	"errors"

	"github.com/stockfolioofficial/back-editfolio/util/pointer"
)

var (
	ItemNotFound = errors.New("item not found")

	UserWrongPassword = errors.New("wrong password")

	UserNotAdmin = errors.New("not admin")

	ItemAlreadyExist = errors.New("item already exsits")

	InvalidateTokenResponse = ErrorResponse{
		ErrorCode: pointer.String("A-1"),
		Message:   "unauthorized",
	}

	NoPermissionResponse = ErrorResponse{
		ErrorCode: pointer.String("A-2"),
		Message:   "no permission",
	}

	UserSignInFailedResponse = ErrorResponse{
		ErrorCode: pointer.String("U-1"),
		Message:   "unauthorized",
	}

	UserWrongPasswordToUpdatePassword = ErrorResponse{
		ErrorCode: pointer.String("U-2"),
		Message:   UserWrongPassword.Error(),
	}

	ItemExist = ErrorResponse{
		ErrorCode: pointer.String("U-3"),
		Message:   ItemAlreadyExist.Error(),
	}

	EmailExistsResponse = ErrorResponse{
		ErrorCode: pointer.String("U-4"),
		Message:   "email exists",
	}

	ServerInternalErrorResponse = ErrorResponse{
		Message: "server internal error",
	}
)

type ErrorResponse struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	Message   string  `json:"message"`
}
