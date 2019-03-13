package tools

import (
	"errors"
)

const (
	ErrMessageSystem = "system error"
)

var (
	ErrUserNameEmpty           = errors.New("the username is empty")
	ErrPasswordEmpty           = errors.New("the password is empty")
	ErrValidateUserPassword    = errors.New("the username and password is mismatching")
	ErrServiceAccountEmpty     = errors.New("the serviceAccount token is empty")
	ErrServiceAccountNotExists = errors.New("the serviceAccount token is not exists")
	ErrParamTidEmpty           = errors.New("the param tid is empty")

	ErrK8sClientInitFailed = errors.New("kubernetes client init failed")
)
