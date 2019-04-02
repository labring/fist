package rbac

import (
	"testing"

	"github.com/wonderivan/logger"
)

func Test_generatorCookieValue(t *testing.T) {
	logger.Info("result:", generatorCookieValue(NewAdminUserInfo("admin", "admin123")))
}
