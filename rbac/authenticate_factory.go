package rbac

//DoFactoryAuthentication is user login access function
func DoFactoryAuthentication(user, password string) *UserInfo {
	var authenticators = map[int]func(user, password string) *UserInfo{1: adminAuth, 2: userAuth, 3: ldapAuth}
	var userInfo *UserInfo
	for _, v := range authenticators {
		userInfo = v(user, password)
		if userInfo != nil {
			return userInfo
		}
	}
	return nil
}

func adminAuth(user, password string) *UserInfo {
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

func userAuth(user, password string) *UserInfo {
	userInfo := GetUserInfo(user, true)
	if userInfo != nil && password == userInfo.Password {
		return userInfo
	}
	return nil
}

func ldapAuth(user, password string) *UserInfo {
	if RbacLdapEnable {
		//if user enable ldap
		if err := authenticationLdap(user, password); err != nil {
			log.Fatal(err)
		    return nil 
		} 
		return NewLdapUserInfo(user, getLdapUserCn(user, password), password )
    }
	return nil
}
