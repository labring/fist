package template

import (
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
)

//Serve start a auth server
func Serve(cmd *cobra.Command) {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	auth := new(restful.WebService)
	//registry  fist auth
	FistRegister(auth)
	wsContainer.Add(auth)
	//cors
	tools.Cors(wsContainer)

	//process port for command
	port, _ := cmd.Flags().GetUint16("port")
	sPort := ":" + strconv.FormatUint(uint64(port), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
