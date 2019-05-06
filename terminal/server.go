package terminal

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	//TerminalPort is cmd port param
	TerminalPort uint16
)

//Serve start a terminal server
func Serve() {
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

func cleanTerminal(namespace string) {
	clientSet := tools.GetK8sClient()
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	t := time.NewTicker(600 * time.Second) //every 10min check heartbeat
	defer t.Stop()
	for {
		select {
		case <-t.C:
			logger.Info("timer running for cleanTerminal.")
			list, err := deploymentsClient.List(metav1.ListOptions{})
			if err != nil {
				log.Fatal(err)
			}
			for _, d := range list.Items {
				var hbInterface Heartbeater
				hbInterface = NewHeartbeater(d.Name, namespace)
				err := hbInterface.CleanTerminalJob()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
