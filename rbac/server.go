package rbac

import (
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"
)

var (
	//RbacPort is cmd port
	RbacPort uint16
	//RbacLdapEnable is cmd enable for ldap
	RbacLdapEnable bool
)

var (
	//RbacLdapPort is config port for ldap . type string
	RbacLdapPort uint16
	//RbacLdapHost is cmd host for ldap
	RbacLdapHost string
	//RbacLdapBindDN is cmd bind-dn for ldap
	RbacLdapBindDN string
	//RbacLdapBindPassword is cmd bind-password for ldap
	RbacLdapBindPassword string
)

//Serve start a auth server
func Serve() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	auth := new(restful.WebService)
	//registry  fist auth
	FistRegister(auth)
	wsContainer.Add(auth)
	//cors
	tools.Cors(wsContainer)
	//process port for command
	sPort := ":" + strconv.FormatUint(uint64(RbacPort), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
