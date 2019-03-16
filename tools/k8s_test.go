package tools

import (
	"fmt"
	"testing"
)

func TestCreateNamespace(t *testing.T) {
	CreateNamespace("ffff")
}

func TestGetSecrets(t *testing.T) {
	secrect := priGetSecrets("sealyun", "fist-admin")
	fmt.Println(string(secrect.Data["password"]))
	fmt.Println(string(secrect.Data["username"]))
}

func TestCreateSecrets(t *testing.T) {
	data := make(map[string]string)
	data["username"] = "cuisongliu"
	data["password"] = "admin"
	data["groups"] = "users"
	data["nickname"] = "cuisongliu1"
	err := SealyunCreateSecretsForMap(UserOperator, "cuisongliu1", data, data)
	if err != nil {
		panic(err)
	}
}

func TestGetSecretsMap(t *testing.T) {
	data := SealyunGetSecretMap(UserOperator, "cuisongliu")
	for index, value := range data {
		println(index, ":", value)
	}
}

func TestListSecrets(t *testing.T) {
	//client := instanceSingleK8sClient()
	//label := strings.Join([]string{"nickname=cuisongliu1"},",")
	//list ,_ := client.CoreV1().Secrets(DefaultNamespace).List(metav1.ListOptions{
	//	LabelSelector:label,
	//})
	//items :=list.Items
	//for _,i := range items{
	//	println(i.Name)
	//}
	data := make(map[string]string)
	data["groups"] = "users"
	arr := priListSecretArr(DefaultNamespace, data)
	for _, v := range arr {
		for i, iv := range v {
			println(i, ":", iv)
		}
		println("next")
	}
}

func TestDeleteSecrets(t *testing.T) {
	err := SealyunDeleteSecrets(UserOperator, "cuisongliu")
	if err != nil {
		panic(err)
	}
}
