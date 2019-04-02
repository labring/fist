package tools

import (
	"errors"
)

//const is global const
const (
	ErrMessageSystem = "system error"
)

//vars is global var
var (
	ErrUserNameEmpty           = errors.New("the username is empty")
	ErrPasswordEmpty           = errors.New("the password is empty")
	ErrValidateUserPassword    = errors.New("the username and password is mismatching")
	ErrServiceAccountEmpty     = errors.New("the serviceAccount token is empty")
	ErrServiceAccountNotExists = errors.New("the serviceAccount token is not exists")
	ErrParamTidEmpty           = errors.New("the param tid is empty")
	ErrUserAuth                = errors.New("user auth failed ")
	ErrUserAdd                 = errors.New("user add failed ")
	ErrUserUpdate              = errors.New("user update failed ")
	ErrUserDel                 = errors.New("user del failed ")
	ErrUserGet                 = errors.New("user get failed ")
	ErrK8sClientInitFailed     = errors.New("kubernetes client init failed")
	//validate
	ErrUserName      = errors.New("user name is error ")
	ErrUserNotExists = errors.New("user not exists ")
	ErrUserExists    = errors.New("user is exists ")
	//
	ErrUserNotLogin = errors.New("user not login ")
	//ldap
	ErrLdapUserNotExists = errors.New("ldap user not exists ")
	//k8s
	ErrAuthTokeKeyError = errors.New("auth loadKeyPair is error ")
)
