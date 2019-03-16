package rbac

import (
	"github.com/emicklei/go-restful"
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
	//registry  fist auth
	FistRegister(auth)
	wsContainer.Add(auth)
	//process port for command
	port, _ := cmd.Flags().GetUint16("port")
	sPort := ":" + strconv.FormatUint(uint64(port), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
