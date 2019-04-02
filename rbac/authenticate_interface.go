package rbac

//DoInterfaceAuthentication is user login access function
//Deprecated
func DoInterfaceAuthentication(user, password string) *UserInfo {
	var authenticators = []authenticator{newAdminAuth(), newLdapAuth(), newKubeSecretAuth()}
	var userInfo *UserInfo
	for _, v := range authenticators {
		userInfo = v.Authenticate(user, password)
		if userInfo != nil {
			return userInfo
		}
	}
	return nil
}

//newAdminAuth construction method  for admin auth
//Deprecated
func newAdminAuth() authenticator {
	var iAuthenticator authenticator
	iAuthenticator = &AdminAuth{}
	return iAuthenticator
}

//newKubeSecretAuth construction method  for user name
//Deprecated
func newKubeSecretAuth() authenticator {
	var iAuthenticator authenticator
	iAuthenticator = &KubeSecretAuth{}
	return iAuthenticator
}

//newLdapAuth construction method  for ldap
//Deprecated
func newLdapAuth() authenticator {
	var iAuthenticator authenticator
	iAuthenticator = &LdapAuth{}
	return iAuthenticator
}

//authenticator interface for auth
type authenticator interface {
	Authenticate(user, password string) *UserInfo //error return nil ,other return userInfo
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
	}
	return nil
}
