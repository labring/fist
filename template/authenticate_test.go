package rbac

import (
	"github.com/wonderivan/logger"
	"testing"
)

func TestDoAuthentication(t *testing.T) {
	userInfo := DoAuthentication("admin", "1f2d1e2e67df")
	if userInfo != nil {
		logger.Info(userInfo.Username, userInfo.Password, userInfo.Nickname, userInfo.Groups)
	} else {
		logger.Error("user not exists")
	}
}
