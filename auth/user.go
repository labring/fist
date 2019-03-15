package auth

//UserInfo is user info struct for restful http
type UserInfo struct {
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

//NewUserInfo construction method
func NewUserInfo(name, password string, groups []string) *UserInfo {
	return &UserInfo{Name: name, Password: password, Groups: groups}
}

//NewAdminUserInfo construction method for admin
func NewAdminUserInfo(name, password string) *UserInfo {
	return &UserInfo{Name: name, Password: password, Groups: []string{"admin"}}
}

//NewLdapUserInfo construction method for ldap
func NewLdapUserInfo(name, password string) *UserInfo {
	return &UserInfo{Name: name, Password: password, Groups: []string{"ldap"}}
}

//func GetUserInfo(name string) *UserInfo  {
//	return nil
//}
//
//func ListAllUserInfo() *[]UserInfo  {
//	return nil
//}
//
//
//func AddUserInfo(userInfo *UserInfo) error {
//	return nil
//}
//
//func UpdateUserInfo(userInfo *UserInfo) error {
//	return nil
//}
//
//func DelUserInfo(name string) error {
//	return nil
//}
