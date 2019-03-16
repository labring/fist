package rbac

import (
	"github.com/fanux/fist/tools"
	"strings"
)

//UserInfo is user info struct for restful http
type UserInfo struct {
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Password string   `json:"password,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

//NewUserInfo is construction method
func NewUserInfo(username, nickname, password string, groups []string) *UserInfo {
	return &UserInfo{Username: username, Nickname: nickname, Password: password, Groups: groups}
}

//NewDefaultUserInfo is construction method
func NewDefaultUserInfo(username, password string, groups []string) *UserInfo {
	return NewUserInfo(username, username, password, groups)
}

//NewAdminUserInfo  is construction method for admin
func NewAdminUserInfo(username, password string) *UserInfo {
	return NewUserInfo(username, "administrator", password, []string{"admin"})
}

//NewLdapUserInfo is construction method for ldap
func NewLdapUserInfo(username, nickname, password string) *UserInfo {
	return NewUserInfo(username, nickname, password, []string{"ldap"})
}

//GetUserInfo is add method from server_fist
func GetUserInfo(name string) *UserInfo {
	secretMap := tools.SealyunGetSecretMap(tools.UserOperator, name)
	if secretMap != nil {
		return mapToUserInfo(secretMap, false)
	}
	return nil
}

//ListAllUserInfo is list method from server_fist
func ListAllUserInfo() []*UserInfo {
	labelsMap := make(map[string]string)
	labelsMap["module"] = "users"
	secrets := tools.SealyunListSecrets(tools.UserOperator, labelsMap)
	var users = make([]*UserInfo, len(secrets))
	if len(secrets) != 0 {
		for i, v := range secrets {
			users[i] = mapToUserInfo(v, false)
		}
	}
	return users
}

//AddUserInfo is add method from server_fist
func AddUserInfo(userInfo *UserInfo) error {
	labelsMap := make(map[string]string)
	labelsMap["module"] = "users"
	err := tools.SealyunCreateSecretsForMap(tools.UserOperator, userInfo.Username, userInfoToMap(userInfo), labelsMap)
	if err != nil {
		return err
	}
	return nil
}

//UpdateUserInfo is update method from server_fist
func UpdateUserInfo(userInfo *UserInfo) error {
	labelsMap := make(map[string]string)
	labelsMap["module"] = "users"
	err := tools.SealyunUpdateSecretsForMap(tools.UserOperator, userInfo.Username, userInfoToMap(userInfo), labelsMap)
	if err != nil {
		return err
	}
	return nil
}

//DelUserInfo is del method from server_fist
func DelUserInfo(name string) error {
	err := tools.SealyunDeleteSecrets(tools.UserOperator, name)
	if err != nil {
		return err
	}
	return nil
}

func mapToUserInfo(data map[string]string, needPwd bool) *UserInfo {
	username := data["username"]
	password := data["password"]
	nickname := data["nickname"]
	groupsStr := data["groups"]
	if tools.NotEmptyAll(username, password) {
		groups := strings.Split(groupsStr, ",")
		if !needPwd {
			password = ""
		}
		return NewUserInfo(username, nickname, password, groups)
	}
	return nil
}
func userInfoToMap(userInfo *UserInfo) map[string]string {
	secretMap := make(map[string]string)
	secretMap["username"] = userInfo.Username
	secretMap["password"] = userInfo.Password
	secretMap["nickname"] = userInfo.Username
	secretMap["groups"] = strings.Join(userInfo.Groups, ",")
	return secretMap
}
