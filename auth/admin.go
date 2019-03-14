package auth

import (
	"errors"
	"github.com/fanux/fist/tools"
)

//AdminUsername is global vars for admin username
var AdminUsername string

//AdminPassword is global vars for admin password
var AdminPassword string

// consts
const (
	DefaultNamespace  = "sealyun"
	DefaultSecretName = "fist-admin"
)

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
	clients, err := tools.GetK8sClient()
	if err != nil {
		return err
	}
	if AdminUsername == "" {
		secrets, err := tools.GetSecrets(DefaultNamespace, DefaultSecretName, clients)
		if err != nil {
			return err
		}
		AdminUsername = string(secrets.Data["username"])
	}
	if AdminPassword == "" {
		secrets, err := tools.GetSecrets(DefaultNamespace, DefaultSecretName, clients)
		if err != nil {
			return err
		}
		AdminPassword = string(secrets.Data["password"])
	}
	return nil
}

// IsAdmin is implements for Adminer  function
func (admin *Admin) IsAdmin() (bool, error) {
	if admin.Name == "" {
		return false, errors.New("the username is empty")
	}
	if admin.Passwd == "" {
		return false, errors.New("the password is empty")
	}
	if admin.Name != AdminUsername || admin.Passwd != AdminPassword {
		return false, errors.New("the username and password is mismatching")
	}
	return true, nil
}
