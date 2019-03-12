package tools

import "github.com/emicklei/go-restful"

type responseObject struct {
	Message string      `json:"message"`
	Code    int32       `json:"code"`
	Data    interface{} `json:"data"`
}

//ResponseSuccess
func ResponseSuccess(response *restful.Response, data interface{}) {
	response.WriteEntity(responseObject{Code: 200, Message: "success", Data: data})
}

//ResponseError
func ResponseErrorAndCode(response *restful.Response, code int32, err error) {
	response.WriteEntity(responseObject{Code: code, Message: err.Error(), Data: ""})
}

func ResponseError(response *restful.Response, err error) {
	ResponseErrorAndCode(response, 500, err)
}
