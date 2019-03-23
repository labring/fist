package auth

import (
	"github.com/emicklei/go-restful"
	"github.com/wonderivan/logger"
	"log"
	"net/http"
	"strconv"
)

var (
	//AuthPort is cmd port param
	AuthHttpsPort uint16
	//AuthCert is cmd cert file
	AuthCert string
	//AuthKey is cmd key file
	AuthKey string
)

//Serve start a auth server
func Serve() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	auth := new(restful.WebService)
	//registry k8s auth and fist auth
	K8sRegister(auth)
	wsContainer.Add(auth)
	//process port for command
	sPort := ":" + strconv.FormatUint(uint64(AuthHttpsPort), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	//process cert/key for command
	logger.Info("certFile is :", AuthCert, ";keyFile is:", AuthKey)

	log.Fatal(server.ListenAndServeTLS(AuthCert, AuthKey))
}
