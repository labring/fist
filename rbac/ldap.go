package rbac

import (
	"fmt"
	"gopkg.in/ldap.v3"
	"log"
	"strings"
)

func getLdapSearchResult(user, password string) (*ldap.SearchResult, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", RbacLdapHost, RbacLdapPort))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind(RbacLdapBindDN, RbacLdapBindPassword)
	if err != nil {
		log.Fatal(err)
	}
	ldapdn := strings.Split(RbacLdapBindDN, ",")
	dc := ldapdn[1] + string(',') + ldapdn[2] // dc=sealyun,dc=com
	searchRequest := ldap.NewSearchRequest(
		dc, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=*)(uid=%s))", user),
		[]string{"dn","cn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
		return sr, err
	}
	return sr, nil
}

func authenticationLdap(user, password string) error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", RbacLdapHost, RbacLdapPort))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	sr, err := getLdapSearchResult(user, password)
	if err != nil {
		return err
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		return err
	}
	return nil
}

func getLdapUserCn(user, password string) string {
	sr, err := getLdapSearchResult(user, password)

	if err != nil {
		log.Fatal(err)
	}

	if len(sr.Entries) == 0 {
        log.Fatal("user not found") 
    }
	return sr.Entries[0].GetAttributeValue("cn") //get nickname by ldap cn
}