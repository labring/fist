package auth

import (
	"errors"
	"github.com/fanux/fist/tools"
)

// vars
var AdminUsername string
var AdminPassword string

// consts
const (
	DefaultNamespace  = "sealyun"
	DefaultSecretName = "fist-admin"
)

type Admin struct {
	Name   string
	Passwd string
}

type Adminer interface {
	LoadSecret() error
	IsAdmin() (bool, error)
}

func NewAdmin(name string, passwd string) Adminer {
	var admire Adminer
	admire = &Admin{Name: name, Passwd: passwd}
	return admire
}

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

func (admin *Admin) IsAdmin() (bool, error) {
	if admin.Name == "" {
		return false, errors.New("the username is empty")
	}
	if admin.Passwd == "" {
		return false, errors.New("the password is empty")
	}
	if admin.Name == AdminUsername && admin.Passwd == AdminPassword {
		return true, nil
	} else {
		return false, errors.New("the username and password is mismatching")
	}
}
