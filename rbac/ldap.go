package rbac 

import (
        "strings"
        "gopkg.in/ldap.v3"
        "log"
        "fmt"
)
 
func getLdapSearchResult (user, password string) (ldap.SearchResult,error) {
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
                fmt.Sprintf("(&(objectClass=*)(uid=%s))", user),
                []string{ "dn"},
                nil,
        )
        sr, err := l.Search(searchRequest)
        if err != nil {
            log.Fatal(err)
            return nil,err
        } else {
            return sr,nil 
        }
}

func authenticationLdap(user, password string) error {
        sr, err := getLdapSearchResult(user, password)
        if err != nil {
            return err
        }
        userdn := sr.Entries[0].DN
        // Bind as the user to verify their password
        err = l.Bind(userdn, password)
        if err != nil {
            return err
        } else {
            log.Fatal("user authenticated")
            return nil 
        }
}

func getLdapUserCn(user, password string) string {
    sr, err := getLdapSearchResult(user, password)
    if err != nil {
        log.Fatal(err)
    }
    return sr.Entries[0].GetAttributeValue("cn") //get nickname by ldap cn
}