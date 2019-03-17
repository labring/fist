package rbac

import (
	"github.com/wonderivan/logger"
	"testing"
)

func TestDoAuthentication(t *testing.T) {
	userInfo := DoAuthentication("user", "123")
	if userInfo != nil {
		logger.Info(userInfo.Username, userInfo.Password, userInfo.Nickname, userInfo.Groups)
	} else {
		logger.Error("user not exists")
	}
}
