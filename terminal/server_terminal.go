package terminal

import (
	"log"
	"time"

	"github.com/fanux/fist/rbac"

	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"

	"github.com/emicklei/go-restful"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Register is
func Register(container *restful.Container) {
	terminal := new(restful.WebService)
	terminal.
		Path("/").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	terminal.Route(terminal.POST("/terminal").Filter(rbac.CookieFilter).To(createTerminal))
	terminal.Route(terminal.GET("/heartbeat").Filter(rbac.CookieFilter).To(handleHeartbeat))

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
