package tools

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"github.com/wonderivan/logger"
)

func pcks5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pcks5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//DESEncrypt is des encrypt
func DESEncrypt(origData, key []byte) string {
	if len(key) != 8 {
		logger.Error("key length must is 8.")
		return ""
	}
	block, _ := des.NewCipher(key)
	origData = pcks5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted)
}

//DESDecrypt is des  decrypt
func DESDecrypt(data string, key []byte) string {
	if len(key) != 8 {
		logger.Error("key length must is 8.")
		return ""
	}
	crypted, _ := base64.StdEncoding.DecodeString(data)
	block, _ := des.NewCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pcks5UnPadding(origData)
	return string(origData)
}
