package tools

import (
	"fmt"
	"testing"
)

func TestCreateNamespace(t *testing.T) {
	CreateNamespace("ffff")
}

func TestGetSecrets(t *testing.T) {
	secrect, _ := GetSecrets("sealyun", "fist-admin")
	fmt.Println(string(secrect.Data["password"]))
	fmt.Println(string(secrect.Data["username"]))
}
