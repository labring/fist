package auth

import (
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"
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
	Pub, Priv = CreateKeyPair()
	go httpServer()
	httpsServer()
}

func httpsServer() {
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

func httpServer() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	auth := new(restful.WebService)
	//registry k8s auth and fist auth
	TokenRegister(auth)
	wsContainer.Add(auth)
	//cors
	tools.Cors(wsContainer)
	//process port for command
	httpPort := ":" + strconv.FormatUint(uint64(AuthHTTPPort), 10)
	logger.Info("start listening on localhost", httpPort)
	server := &http.Server{Addr: httpPort, Handler: wsContainer}

	log.Fatal(server.ListenAndServe())
}
