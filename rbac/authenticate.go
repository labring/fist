package rbac

import (
	"strings"
	"gopkg.in/ldap.v3"
 	"log"
 	"fmt"

)
//DoAuthentication is user login access function
func DoAuthentication(user, password string) *UserInfo {
	var authenticators = []authenticator{newAdminAuth(), newLdapAuth(), newKubeSecretAuth()}
	var userInfo *UserInfo
	for i := 0; i < len(authenticators); i++ {
		userInfo = authenticators[i].Authenticate(user, password)
		if userInfo != nil {
			return userInfo
		}
	}
	return nil
}

//newAdminAuth construction method  for admin auth
func newAdminAuth() authenticator {
	var iAuthenticator authenticator
	iAuthenticator = &AdminAuth{}
	return iAuthenticator
}

//newKubeSecretAuth construction method  for user name
func newKubeSecretAuth() authenticator {
	var iAuthenticator authenticator
	iAuthenticator = &KubeSecretAuth{}
	return iAuthenticator
}

//newLdapAuth construction method  for ldap
func newLdapAuth() authenticator {
	var iAuthenticator authenticator
	iAuthenticator = &LdapAuth{}
	return iAuthenticator
}

//authenticator interface for auth
type authenticator interface {
	Authenticate(user, password string) *UserInfo //失败返回nil  成功返回用户信息
}

//AdminAuth is struct type
type AdminAuth struct{}

//KubeSecretAuth is struct type
type KubeSecretAuth struct{}

//LdapAuth is struct type
type LdapAuth struct{}

//Authenticate is interface impl for AdminAuth
func (AdminAuth) Authenticate(user, password string) *UserInfo {
	admire := NewAdmin(user, password)
	err := admire.LoadSecret()
	if err != nil {
		return nil
	}
	isAdmin, err := admire.IsAdmin()
	if err == nil && isAdmin {
		return NewAdminUserInfo(user, "")
	}
	return nil
}

//Authenticate is interface impl for KubeSecretAuth
func (KubeSecretAuth) Authenticate(user, password string) *UserInfo {
	userInfo := GetUserInfo(user, true)
	if userInfo != nil && password == userInfo.Password {
		return userInfo
	}
	return nil
}

//Authenticate is interface impl for LdapAuth
func (LdapAuth) Authenticate(user, password string) *UserInfo {
	if RbacLdapEnable {
		//if user enable ldap
        l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", RbacLdapHost, RbacLdapPort))
        if err != nil {
                log.Fatal(err)
        }
        defer l.Close()

        err = l.Bind(RbacLdapBindDN, RbacLdapBindPassword)
        if err != nil {
                log.Fatal(err)
        }
        ldapdn :=  strings.Split(RbacLdapBindDN, ",")    
        dc := ldapdn[1] + string(',') +  ldapdn[2] // dc=sealyun,dc=com
        searchRequest := ldap.NewSearchRequest(
                dc, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
                fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", user),
                []string{ "dn"},
                nil,
        )
        sr, err := l.Search(searchRequest)
        if err != nil {
                log.Fatal(err)
        }

        userdn := sr.Entries[0].DN
        // Bind as the user to verify their password
        err = l.Bind(userdn, password)
        if err != nil {
                log.Fatal(err)
        }
        log.Fatal("user authenticated")

		return &UserInfo{Username: user , Nickname: sr.Entries[0].GetAttributeValue("cn"), Password: password}
	}
	return nil
}
