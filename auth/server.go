package auth

import (
	"github.com/emicklei/go-restful"
	"github.com/wonderivan/logger"
	"log"
	"net/http"
	"strconv"
)

var (
	//AuthHTTPSPort is cmd port param
	AuthHTTPSPort uint16
	//AuthHTTPPort is cmd port param
	AuthHTTPPort uint16
	//AuthCert is cmd cert file
	AuthCert string
	//AuthKey is cmd key file
	AuthKey string

	//authHTTPSPortString is string of AuthHTTPSPort
	authHTTPSPortString string
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
	authHTTPSPortString = ":" + strconv.FormatUint(uint64(AuthHTTPSPort), 10)
	logger.Info("start listening on localhost", authHTTPSPortString)
	server := &http.Server{Addr: authHTTPSPortString, Handler: wsContainer}
	//process cert/key for command
	logger.Info("certFile is :", AuthCert, ";keyFile is:", AuthKey)

	log.Fatal(server.ListenAndServeTLS(AuthCert, AuthKey))
}
