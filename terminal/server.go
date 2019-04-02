package terminal

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"

	"github.com/emicklei/go-restful"
)

var (
	//TerminalPort is cmd port param
	TerminalPort uint16
)

//Serve start a terminal server
func Serve() {
	LoadTerminalID()
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	Register(wsContainer)
	//cors
	tools.Cors(wsContainer)

	//clean dead terminal
	go cleanTerminal(DefaultTTYnameapace)

	//process port for command
	sPort := ":" + strconv.FormatUint(uint64(TerminalPort), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
