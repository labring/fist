package rbac

import (
	"regexp"

	"github.com/fanux/fist/tools"
)

// validate username  , only  a-z A-Z
func validateUserName(username string) bool {
	matchChar, _ := regexp.Match("[a-zA-Z]", []byte(username))
	return matchChar
}

func validateUserNameExist(username string) bool {
	if username != "admin" {
		data := tools.SealyunGetSecretMap(tools.UserOperator, username)
		if _, ok := data["username"]; ok {
			return true
		}
		return false
	}
	return true
}

func validateGroups(groups []string) bool {
	for _, val := range groups {
		if val == "administrator" || val == "ldap" {
			return false
		}
	}
	return true
}
