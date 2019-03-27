package tools

import (
	"github.com/wonderivan/logger"
	"testing"
)

func TestDESEncrypt(t *testing.T) {
	data := []byte("hello world")
	key := []byte("df9gtsq3")

	str := DESEncrypt(data, key)
	logger.Info("encrypt str ", str)

	logger.Info("back str", DESDecrypt(str, key))
}
