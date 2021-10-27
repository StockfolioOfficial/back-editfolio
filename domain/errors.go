package domain

import (
	"errors"

	"github.com/stockfolioofficial/back-editfolio/util/pointer"
)

var (
	ErrItemNotFound = errors.New("item not found")

	ErrUserWrongPassword = errors.New("wrong password")

	ErrUserNotAdmin = errors.New("not admin")

	ErrNoPermission = errors.New("no permission")

	ErrItemAlreadyExist = errors.New("item already exsits")

	InvalidateTokenResponse = ErrorResponse{
		ErrorCode: pointer.String("A-1"),
		Message:   "unauthorized",
	}

	NoPermissionResponse = ErrorResponse{
		ErrorCode: pointer.String("A-2"),
		Message:   ErrNoPermission.Error(),
	}

	UserSignInFailedResponse = ErrorResponse{
		ErrorCode: pointer.String("U-1"),
		Message:   "unauthorized",
	}

	UserWrongPasswordToUpdatePassword = ErrorResponse{
		ErrorCode: pointer.String("U-2"),
		Message:   ErrUserWrongPassword.Error(),
	}

	ItemExist = ErrorResponse{
		ErrorCode: pointer.String("U-3"),
		Message:   ErrItemAlreadyExist.Error(),
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
