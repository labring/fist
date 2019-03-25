package rbac

import (
	"github.com/wonderivan/logger"
	"testing"
)

func TestDoInterfaceAuthentication(t *testing.T) {
	userInfo := DoInterfaceAuthentication("admin", "1f2d1e2e67df")
	if userInfo != nil {
		logger.Info(userInfo.Username, userInfo.Password, userInfo.Nickname, userInfo.Groups)
	} else {
		logger.Error("user not exists")
	}
}
func TestDoFactoryAuthentication(t *testing.T) {
	userInfo := DoFactoryAuthentication("admin", "1f2d1e2e67df")
	if userInfo != nil {
		logger.Info(userInfo.Username, userInfo.Password, userInfo.Nickname, userInfo.Groups)
	} else {
		logger.Error("user not exists")
	}
}
