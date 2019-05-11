package template

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
)

//Register is
func Register(container *restful.Container) {
	LoadTemplates("")
	template := new(restful.WebService)
	template.
		Path("/").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	template.Route(template.POST("/templates").To(createTemplates))
	container.Add(template)
}

func createTemplates(request *restful.Request, response *restful.Response) {
	t := request.QueryParameter("type")
	tps := new([]Value)
	err := json.NewDecoder(request.Request.Body).Decode(tps)
	if err != nil {
		log.Fatal(err)
	}

	res := new([]string)
	for _, t := range *tps {
		tempres := RenderValue(t)
		if tempres != "" {
			*res = append(*res, tempres)
		}
	}
	if t == "text" {
		response.AddHeader("Content-type", "text/plain")
		var ss string
		for _, s := range *res {
			ss = fmt.Sprintf("%s\n---\n%s", ss, s)
		}

		response.ResponseWriter.Write([]byte(ss))
		return
	}
	response.WriteEntity(res)
	return
}

//Serve start a template server
func Serve(cmd *cobra.Command) {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	Register(wsContainer)
	//cors
	tools.Cors(wsContainer)
	//process port for command
	port, _ := cmd.Flags().GetUint16("port")
	sPort := ":" + strconv.FormatUint(uint64(port), 10)
	logger.Info("start listening on localhost", sPort)
	server := &http.Server{Addr: sPort, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
