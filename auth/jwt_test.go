package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

func TestCreateKeyPairPri(t *testing.T) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA FIST PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return
	}
	err = pem.Encode(file, block)
	if err != nil {
		return
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	block = &pem.Block{
		Type:  "RSA FIST PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return
	}
	err = pem.Encode(file, block)
	if err != nil {
		return
	}

}

func TestLoadKeyPairPri(t *testing.T) {
	//PubKey, _ := ioutil.ReadFile("public.pem")
	block, _ := pem.Decode(publicPem)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	_ = pubInterface.(*rsa.PublicKey)

	//PriKey, _ := ioutil.ReadFile("private.pem")
	block, _ = pem.Decode(privatePem)
	if block == nil {
		return
	}
	_, err = x509.ParsePKCS1PrivateKey(block.Bytes)
}
