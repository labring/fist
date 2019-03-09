package auth

import (
	"errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// vars
var AdminUsername string
var AdminPassword string

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
	clients, err := GetK8sClient()
	if err != nil {
		return err
	}
	if AdminUsername == "" {
		secrets, err := GetSecrets("sealyun", "fist-admin", clients)
		if err != nil {
			return err
		}
		AdminUsername = string(secrets.Data["username"])
	}
	if AdminPassword == "" {
		secrets, err := GetSecrets("sealyun", "fist-admin", clients)
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

func GetK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func GetSecrets(namespace string, name string, clientset *kubernetes.Clientset) (*v1.Secret, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	return secret, err
}
