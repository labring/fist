package tools

import (
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//operator k8s

//priGetSecrets is get secrets for k8s
func priGetSecrets(namespace, name string) *v1.Secret {
	client := instanceSingleK8sClient()
	secret, err := client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		logger.Error("priGetSecrets is  error: ", err)
		return nil
	}
	return secret
}

//priDeleteSecrets is del secrets for k8s
func priListSecrets(namespace string, labels map[string]string) []v1.Secret {
	client := instanceSingleK8sClient()

	secretsList, err := client.CoreV1().Secrets(namespace).List(metav1.ListOptions{
		LabelSelector: MapToString(labels),
	})
	if err != nil {
		logger.Error("priListSecrets is  error: ", err)
		return make([]v1.Secret, 0)
	}
	return secretsList.Items
}

//priDeleteSecrets is del secrets for k8s
func priDeleteSecrets(namespace, name string) error {
	client := instanceSingleK8sClient()
	err := client.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		logger.Error("priGetSecrets is  error: ", err)
		return err
	}
	return nil
}

//priCreateSecrets is create secrets for k8s
func priCreateSecrets(namespace string, secret *v1.Secret) error {
	client := instanceSingleK8sClient()
	_, err := client.CoreV1().Secrets(namespace).Create(secret)
	if err != nil {
		logger.Error("priCreateSecrets is  error: ", err)
		return err
	}
	return nil
}

//priUpdateSecrets is update secrets for k8s
func priUpdateSecrets(namespace string, secret *v1.Secret) error {
	client := instanceSingleK8sClient()
	_, err := client.CoreV1().Secrets(namespace).Update(secret)
	if err != nil {
		logger.Error("priUpdateSecrets is  error: ", err)
		return err
	}
	return nil
}

// OperatorType is enum of operator type
type OperatorType int

// consts
const (
	DefaultNamespace       = "sealyun"
	DefaultAdminSecretName = "fist-admin"

	UserOperator OperatorType = 1 + iota

	UserSecretPrefix = "fist-users-"
)

func priGetPrefix(operator OperatorType) string {
	var returnStr string
	switch operator {
	case UserOperator:
		{
			returnStr = UserSecretPrefix
			break
		}
	}
	return returnStr
}

//SealyunGetSecretMap is Sealyun function get map
func SealyunGetSecretMap(operatorType OperatorType, name string) map[string]string {
	return priGetSecretMap(DefaultNamespace, priGetPrefix(operatorType)+name)
}

//SealyunGetSecretString is Sealyun function get str
func SealyunGetSecretString(operatorType OperatorType, name, key string) string {
	return priGetSecretString(DefaultNamespace, priGetPrefix(operatorType)+name, key)
}

//SealyunCreateSecretsForMap is create function for map
func SealyunCreateSecretsForMap(operatorType OperatorType, name string, data, labels map[string]string) error {
	return priCreateSecretsForMap(DefaultNamespace, priGetPrefix(operatorType)+name, data, labels)
}

//SealyunUpdateSecretsForMap is update function for map
func SealyunUpdateSecretsForMap(operatorType OperatorType, name string, data, labels map[string]string) error {
	return priUpdateSecretsForMap(DefaultNamespace, priGetPrefix(operatorType)+name, data, labels)
}

//SealyunUpdateSecretsForString is update function for map
func SealyunUpdateSecretsForString(operatorType OperatorType, name, key, value string, labels map[string]string) error {
	data := SealyunGetSecretMap(operatorType, name)
	data[key] = value
	return priUpdateSecretsForMap(DefaultNamespace, priGetPrefix(operatorType)+name, data, labels)
}

//SealyunDeleteSecrets is delete function
func SealyunDeleteSecrets(operatorType OperatorType, name string) error {
	return priDeleteSecrets(DefaultNamespace, priGetPrefix(operatorType)+name)
}

//SealyunListSecrets is list function
func SealyunListSecrets(operatorType OperatorType, labels map[string]string) []map[string]string {
	return priListSecretArr(DefaultNamespace, labels)
}

//SealyunGetAdminSecretString is Sealyun  admin module function
func SealyunGetAdminSecretString(key string) string {
	return priGetSecretString(DefaultNamespace, DefaultAdminSecretName, key)
}

//public function  common

// priCreateSecretsForMap is create function for map
func priCreateSecretsForMap(namespace, name string, data, labels map[string]string) error {
	object := metav1.ObjectMeta{Name: name, Labels: labels}
	secret := &v1.Secret{StringData: data, ObjectMeta: object}
	err := priCreateSecrets(namespace, secret)
	if err != nil {
		return err
	}
	return nil
}

// priCreateSecretsForMap is create function for map
func priUpdateSecretsForMap(namespace, name string, data, labels map[string]string) error {
	object := metav1.ObjectMeta{Name: name, Labels: labels}
	secret := &v1.Secret{StringData: data, ObjectMeta: object}
	err := priUpdateSecrets(namespace, secret)
	if err != nil {
		return err
	}
	return nil
}

//priGetSecretMap is get function for map
func priListSecretArr(namespace string, labels map[string]string) []map[string]string {
	secretsArr := priListSecrets(namespace, labels)
	returnArr := make([]map[string]string, len(secretsArr))
	if len(secretsArr) != 0 {
		secrets := priListSecrets(namespace, labels)
		for i, secret := range secrets {
			returnArr[i] = priSecretsToMap(&secret)
		}
	}
	return returnArr
}

//priGetSecretMap is get function for map
func priGetSecretMap(namespace, name string) map[string]string {
	secrets := priGetSecrets(namespace, name)
	return priSecretsToMap(secrets)
}

//priGetSecretString is get function for string
func priGetSecretString(namespace, name, key string) string {
	secrets := priGetSecrets(namespace, name)
	if secrets != nil {
		return string(secrets.Data[key])
	}
	return ""
}

//priSecretsToMap is  Secrets convert map
func priSecretsToMap(secret *v1.Secret) map[string]string {
	if secret != nil {
		returnMap := make(map[string]string)
		for index, value := range secret.Data {
			returnMap[index] = string(value)
		}
		return returnMap
	}
	return make(map[string]string, 0)
}
