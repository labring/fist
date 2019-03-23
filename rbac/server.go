package rbac

import (
	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"
	"log"
	"net/http"
	"strconv"
)

var (
	//RbacPort is cmd port
	RbacPort uint16
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
