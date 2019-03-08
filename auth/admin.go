package main

import (
	"errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// vars
var ADMIN_USERNAME string
var ADMIN_PASSWORD string

type Admin struct {
	Name   string
	Passwd string
}

type Adminer interface {
	LoadSecret() error
	IsAdmin() (bool, error)
}

func newAdmin(name string, passwd string) Adminer {
	var adminer Adminer
	adminer = &Admin{Name: name, Passwd: passwd}
	return adminer
}

func (*Admin) LoadSecret() error {
	clientset, err := GetK8sClient()
	if err != nil {
		return err
	}
	if ADMIN_USERNAME == "" {
		secrets, err := GetSecrets("sealyun", "fist-admin", clientset)
		if err != nil {
			return err
		}
		ADMIN_USERNAME = string(secrets.Data["username"])
	}
	if ADMIN_PASSWORD == "" {
		secrets, err := GetSecrets("sealyun", "fist-admin", clientset)
		if err != nil {
			return err
		}
		ADMIN_PASSWORD = string(secrets.Data["password"])
	}
	return nil
}

func (admin *Admin) IsAdmin() (bool, error) {
	if admin.Name == "" {
		return false, errors.New("the username is empty.")
	}
	if admin.Passwd == "" {
		return false, errors.New("the password is empty.")
	}
	if admin.Name == ADMIN_USERNAME && admin.Passwd == ADMIN_PASSWORD {
		return true, nil
	} else {
		return false, errors.New("the username and password is mismatching.")
	}
}

func GetK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func GetSecrets(namespace string, name string, clientset *kubernetes.Clientset) (*v1.Secret, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	return secret, err
}
