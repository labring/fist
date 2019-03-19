package rbac

import (
	"github.com/fanux/fist/tools"
	"regexp"
)

// validate username  , only  a-z A-Z
func validateUserName(username string) bool {
	matchChar, _ := regexp.Match("[a-zA-Z]", []byte(username))
	return matchChar
}

func validateUserNameExist(username string) bool {
	data := tools.SealyunGetSecretMap(tools.UserOperator, username)
	if _, ok := data["username"]; ok {
		return true
	}
	return false
}

func validateGroups(groups []string) bool {
	for _, val := range groups {
		if val == "administrator" || val == "ldap" {
			return false
		}
	}
	return true
}
