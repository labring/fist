package tools

import (
	"errors"
)

var (
	ErrUserNameEmpty           = errors.New("the username is empty")
	ErrPasswordEmpty           = errors.New("the password is empty")
	ErrValidateUserPassword    = errors.New("the username and password is mismatching")
	ErrServiceAccountEmpty     = errors.New("the serviceAccount token is empty")
	ErrServiceAccountNotExists = errors.New("the serviceAccount token is not exists")
	ErrParamTidEmpty           = errors.New("the param tid is empty")

	ErrSignPayload     = errors.New("failed to sign payload")
	ErrSerializeClaims = errors.New("could not serialize claims")
)
