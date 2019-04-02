package tools

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

var publicPem = []byte(`-----BEGIN RSA FIST PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0W1PGq93sl8ehtcf2X4y
UkNIugjTjhRjGfI7W/KkXSub5jUAnU6OMB5cntOisINHQTypzMqNMUvfZUp7c6zj
RgGoWfSH284XZbpEJpjUjEBQSy1YpHUXOGfcQ+Cgvloh43w0MbHzauRyk9usE+Li
cYTMA4lMaY8kZjD8b+gSbbrJpGJz/oBLx0kXJaaoB0ZI4+XytWu5AMZFPdfDAZOy
6OC3JzcLXmywCddPJ/VxTLDDxAEPkJLpD1WnRhGDpY8rp7LnCSbqWFSULDaXcOsw
VAf5EvGqL7YOXvXn+1A3h0Pg/jFkLeGnMtchYebMsKIFpklPPRwMtu5Bd5jZaFoo
PwIDAQAB
-----END RSA FIST PUBLIC KEY-----`)
var privatePem = []byte(`-----BEGIN RSA FIST PRIVATE KEY-----
MIIEpAIBAAKCAQEA0W1PGq93sl8ehtcf2X4yUkNIugjTjhRjGfI7W/KkXSub5jUA
nU6OMB5cntOisINHQTypzMqNMUvfZUp7c6zjRgGoWfSH284XZbpEJpjUjEBQSy1Y
pHUXOGfcQ+Cgvloh43w0MbHzauRyk9usE+LicYTMA4lMaY8kZjD8b+gSbbrJpGJz
/oBLx0kXJaaoB0ZI4+XytWu5AMZFPdfDAZOy6OC3JzcLXmywCddPJ/VxTLDDxAEP
kJLpD1WnRhGDpY8rp7LnCSbqWFSULDaXcOswVAf5EvGqL7YOXvXn+1A3h0Pg/jFk
LeGnMtchYebMsKIFpklPPRwMtu5Bd5jZaFooPwIDAQABAoIBACBR4D1cCvrqFwn0
NSQZh6aOX7DmH8/zcX/VlDeScK84Raz3TQr5zpO6e8y/GIJpIBv/Qq0qLTQIY4d8
QSUO71wkcVjAlh7j5VR7tHzIZTTnz/xqGR34PAcmcCXBis0Vl9lFl8B7l/dNHimX
Yy0GfK/MMLfy/mnj/1gOA0c5524rd8EOauAU68MpSsBkK/Iv8uvjEB7sa564CV7n
vbu/HKvXisvl+ekUONh/Zol6udOyBl7zVDRorOS7FXwMLivqv7HRaDoDvzXxAD8h
MYUZDxfsq54rajvSiQC97I2kNYVy+H4ZwreNa2IPoSz/xi7DmRUFOr9oXGyxypNW
bdXVDdECgYEA16n/gr6E8suAIXJg+nKaxdxsfjXIpF41rgzUytAQiR+OKlkQGBQx
6zoe8iZf4vh153ZaeB2Z8J+R/x/TO5HDhOJUtJ1/MPMQ9fFl/wxf3z1OgTx4//q1
QLejwK5xR2xGYK9Od5OFUQ79329+pTtGh/amqvD8MQ2pgf8GYh3I5iMCgYEA+Jiu
EYjiSw9mZjE/zBjmTr1p3O8YTzo9XUFItjGxhfJSVNScwTJ77CCPB2FumAm3fhHF
fjHmW9/+7ki52LRZGgZj94KWiDKvvmwvswNVZgWRvaAgeDX/ArmFjbitpwO965oB
IBqVxRFlc6SGWJmB99/FEjvymv0yIx5wHS0IITUCgYA9BsTfyWuzjLRQQp6AoEd9
r9cfi7agFGyaONVKIsBbHQvMnfE47xYFGyhAm21Mu8QZYFtPXAkAlxboG6hZVfD+
vFD93wdr4YwuHEYheu3yLNX3KhaPUPFTQ+PZlUNG07PCjjjlC2CRRG3AeTtcJD7c
IIafHpwugVAzEehLgWQlAwKBgQDtKOtXOxl3tjdgYreURGTi1Xz7EvZEDzGfl7qA
ZiQPvz+gQj1RapIikGUjC2ZwkUSGdvnMDFrRZ35TYPo7rMIcO+B+rgxh3skehyNy
SkncTI/fMbYIOKsRtF+e9oJkBQPYklFkiTg8iv4YNusb90awbMAbRymJhuef7VfT
3bQVyQKBgQDSfFa/psuKEdm+Lai3S+3WBwgTITl6y2dCLZ3CftDSTv+7XdLU9bi0
5KkcCgok/uObWGSZZFCaMDEZ3B8EiUqcIdPfMpjNi9PX8ZsxLb2o8t3rdh+Oe2Kc
hrl7TRewrgKHnJt57QXEz5X8/Ov2+C9h2TULunbwTh5iN1Bj1eoo2Q==
-----END RSA FIST PRIVATE KEY-----`)

func PemDefaultPrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode(privatePem)
	if block == nil {
		return nil
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}
	return private
}

func PemDefaultPublicKey() *rsa.PublicKey {
	block, _ := pem.Decode(publicPem)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}
	switch pubInterface := pubInterface.(type) {
	case *rsa.PublicKey:
		return pubInterface
	default:
		break // fall through
	}
	return nil
}

//ParseRsaPubKeyFromPemFile is read publicKey from file
func ParseRsaPubKeyFromPemFile(pemFile string) (*rsa.PublicKey, error) {
	fileIO, err := os.Open(pemFile)
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
	return nil, errors.New("key type is not RSA type")
}

//ParseRsaPrivateKeyFromPemFile is read privateKey from file
func ParseRsaPrivateKeyFromPemFile(pemFile string) (*rsa.PrivateKey, error) {
	fileIO, err := os.Open(pemFile)
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

//ExportRsaPrivateKeyAsPemStr is store privateKey into file
func ExportRsaPrivateKeyAsPemStr(filename string, privkey *rsa.PrivateKey) error {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	file, err := os.Create(filename)
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
