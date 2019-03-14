package tools

import "github.com/emicklei/go-restful"

type responseObject struct {
	Message string      `json:"message"`
	Code    int32       `json:"code"`
	Data    interface{} `json:"data"`
}

//ResponseSuccess is web response for success
func ResponseSuccess(response *restful.Response, data interface{}) {
	response.WriteEntity(responseObject{Code: 200, Message: "success", Data: data})
}

//ResponseError is web response for success
func ResponseError(response *restful.Response, err error) {
	response.WriteEntity(responseObject{Code: 500, Message: err.Error(), Data: ""})
}
