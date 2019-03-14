package tools

import (
	"github.com/wonderivan/logger"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//GetSecrets is get secrets for k8s
func priGetSecrets(namespace, name string) *v1.Secret {
	client := instanceSingleK8sClient()
	secret, err := client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		logger.Error("GetSecrets is  error: ", err)
		return nil
	}
	return secret
}

// consts
const (
	DefaultNamespace  = "sealyun"
	DefaultSecretName = "fist-admin"
)

//public sealyun function
func SealyunGetSecretString(name, key string) string {
	return GetSecretString(DefaultNamespace, name, key)
}
func SealyunGetAdminSecretString(key string) string {
	return GetSecretString(DefaultNamespace, DefaultSecretName, key)
}

//public function  common

func SetSecretsForString(namespace, name, key, value string) {

}

func SetSecretsForMap(namespace, name string, data map[string]string) {

}

func GetSecretMap(namespace, name string) map[string][]byte {
	secrets := priGetSecrets(namespace, name)
	if secrets != nil {
		return secrets.Data
	}
	return nil
}

func GetSecretString(namespace, name, key string) string {
	secrets := priGetSecrets(namespace, name)
	if secrets != nil {
		return string(secrets.Data[key])
	}
	return ""
}
