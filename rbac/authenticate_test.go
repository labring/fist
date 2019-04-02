package rbac

import (
	"testing"
)

func TestDoInterfaceAuthentication(t *testing.T) {
	userInfo := DoInterfaceAuthentication("admin", "1f2d1e2e67df")
	if userInfo != nil {
		t.Log(userInfo.Username, userInfo.Password, userInfo.Nickname, userInfo.Groups)
	} else {
		t.Error("user not exists")
	}
}
func TestDoFactoryAuthentication(t *testing.T) {
	userInfo := DoFactoryAuthentication("admin", "1f2d1e2e67df")
	if userInfo != nil {
		t.Log(userInfo.Username, userInfo.Password, userInfo.Nickname, userInfo.Groups)
	} else {
		t.Error("user not exists")
	}
}

func TestDoFactoryAuthenticationLdap(t *testing.T) {
	RbacLdapEnable = true
	RbacLdapHost = "fist.lameleg.com"
	RbacLdapPort = 31389
	RbacLdapBindDN = "cn=admin,dc=sealyun,dc=com"
	RbacLdapBindPassword = "admin"
	userInfo := DoFactoryAuthentication("fanux", "fanux")
	if userInfo != nil {
		t.Log(userInfo.Username, userInfo.Password, userInfo.Nickname, userInfo.Groups)
	} else {
		t.Error("user not exists")
	}
}
