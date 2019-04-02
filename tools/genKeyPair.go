package tools

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

//ExportRsaPrivateKeyAsPemStr is store privateKey into file
func ExportRsaPrivateKeyAsPemStr(fileDir string, privkey *rsa.PrivateKey) error {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	file, err := os.Create(fileDir)
	if err != nil {
		return err
	}
	err = pem.Encode(file,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

//ParseRsaPrivateKeyFromPemStr is read privateKey from file
func ParseRsaPrivateKeyFromPemStr(fileDir string) (*rsa.PrivateKey, error) {
	fileIO, err := os.Open(fileDir)
	if err != nil {
		return nil, err
	}
	defer fileIO.Close()
	privPEM, err := ioutil.ReadAll(fileIO)
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

//ExportRsaPublicKeyAsPemStr is store publicKey into file
func ExportRsaPublicKeyAsPemStr(fileDir string, pubkey *rsa.PublicKey) error {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return err
	}
	file, err := os.Create(fileDir)
	if err != nil {
		return err
	}
	err = pem.Encode(file,
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

//ParseRsaPublicKeyFromPemStr is read publicKey from file
func ParseRsaPublicKeyFromPemStr(fileDir string) (*rsa.PublicKey, error) {
	fileIO, err := os.Open(fileDir)
	if err != nil {
		return nil, err
	}
	defer fileIO.Close()
	pubPEM, err := ioutil.ReadAll(fileIO)
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}
