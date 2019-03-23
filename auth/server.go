package auth

import (
	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
	"log"
	"net/http"
	"strconv"
)

//Serve start a auth server
func Serve(cmd *cobra.Command) {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	auth := new(restful.WebService)
	//registry k8s auth and fist auth
	K8sRegister(auth)
	wsContainer.Add(auth)
	//cors
	tools.Cors(wsContainer)
	//process port for command
	port, _ := cmd.Flags().GetUint16("port")
	sPort := ":" + strconv.FormatUint(uint64(port), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	//process cert/key for command
	cert, _ := cmd.Flags().GetString("cert")
	key, _ := cmd.Flags().GetString("key")
	logger.Info("certFile is :", cert, ";keyFile is:", key)

	log.Fatal(server.ListenAndServeTLS(cert, key))
}
