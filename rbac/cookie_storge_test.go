package rbac

import (
	"github.com/wonderivan/logger"
	"testing"
)

func Test_generatorCookieValue(t *testing.T) {
	logger.Info("result:", generatorCookieValue(NewAdminUserInfo("admin", "admin123")))
}
