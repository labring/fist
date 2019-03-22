package terminal

import (
	"github.com/fanux/fist/tools"
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

//Register is
func Register(container *restful.Container) {
	terminal := new(restful.WebService)
	terminal.
		Path("/").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	terminal.Route(terminal.POST("/terminal").To(createTerminal))
	terminal.Route(terminal.GET("/heartbeat").To(handleHeartbeat))

	container.Add(terminal)
}

func createTerminal(request *restful.Request, response *restful.Response) {
	t := newTerminal()
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}

	err = t.Create()
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	tools.ResponseSuccess(response, t)
}

func handleHeartbeat(request *restful.Request, response *restful.Response) {
	//get client of k8s
	tid := request.QueryParameter("tid")
	if tid == "" {
		tools.ResponseError(response, tools.ErrParamTidEmpty)
		return
	}
	namespace := request.QueryParameter("namespace")
	if namespace == "" {
		namespace = DefaultTTYnameapace
	}
	var hbInterface Heartbeater
	hbInterface = NewHeartbeater(tid, namespace)
	err := hbInterface.UpdateTimestamp()
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	tools.ResponseSuccess(response, nil)
}

func cleanTerminal(namespace string) {
	clientSet := tools.GetK8sClient()
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)

	go func {
		time.Sleep(time.Duration(600) * time.Second)//every 10min check heartbeat

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

//Serve start a terminal server
func Serve(cmd *cobra.Command) {
	LoadTerminalID()

	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	Register(wsContainer)
	//cors
	tools.Cors(wsContainer)

	//clean dead terminal
	cleanTerminal(DefaultTTYnameapace)

	//process port for command
	port, _ := cmd.Flags().GetUint16("port")
	sPort := ":" + strconv.FormatUint(uint64(port), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}