package auth

import (
	"github.com/fanux/fist/tools"
)

//AdminUsername is global vars for admin username
var AdminUsername string

//AdminPassword is global vars for admin password
var AdminPassword string

// Admin is struct of Adminer interface
type Admin struct {
	Name   string
	Passwd string
}

// Adminer is admin.go interface
type Adminer interface {
	LoadSecret() error
	IsAdmin() (bool, error)
}

// NewAdmin is Admin struct init function
func NewAdmin(name string, passwd string) Adminer {
	var admire Adminer
	admire = &Admin{Name: name, Passwd: passwd}
	return admire
}

// LoadSecret is implements for Adminer  function
func (*Admin) LoadSecret() error {
	if AdminUsername == "" {
		AdminUsername = tools.SealyunGetAdminSecretString("username")
	}
	if AdminPassword == "" {
		AdminPassword = tools.SealyunGetAdminSecretString("password")
	}
	return nil
}

// IsAdmin is implements for Adminer  function
func (admin *Admin) IsAdmin() (bool, error) {
	if admin.Name == "" {
		return false, tools.ErrUserNameEmpty
	}
	if admin.Passwd == "" {
		return false, tools.ErrPasswordEmpty
	}
	if admin.Name != AdminUsername || admin.Passwd != AdminPassword {
		return false, tools.ErrValidateUserPassword
	}
	return true, nil
}
