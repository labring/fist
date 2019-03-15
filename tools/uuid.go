package tools

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

//NewUUID is generator uuid
func NewUUID() string {
	u := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, u); err != nil {
		panic(err)
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u)
}
