package terminal

import (
	"github.com/fanux/fist/rbac"

	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
)

//Register is
func Register(container *restful.Container) {
	terminal := new(restful.WebService)
	var filter restful.FilterFunction
	if RbacEnable {
		filter = rbac.CookieFilter
	} else {
		filter = func(*restful.Request, *restful.Response, *restful.FilterChain) {}
	}
	terminal.
		Path("/").
		Filter(filter).
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
	var hbInterface Heartbeater
	hbInterface = NewHeartbeater(tid, DefaultTTYnameapace)
	err := hbInterface.UpdateTimestamp()
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	tools.ResponseSuccess(response, nil)
}
