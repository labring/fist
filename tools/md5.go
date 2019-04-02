package tools

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5 is method for md5
func MD5(data string) string {
	// md5
	h := md5.New()
	h.Write([]byte(data))
	md5Value := hex.EncodeToString(h.Sum(nil))
	return md5Value
}
