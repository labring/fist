package tools

import (
	"fmt"
	"github.com/emicklei/go-restful"
)

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
func ResponseErrorAndCodeMessage(response *restful.Response, code int32, err error, message string) {
	fmt.Printf("response error: %v", err)
	response.WriteEntity(responseObject{Code: code, Message: message, Data: ""})
}

//custom error
func ResponseError(response *restful.Response, err error) {
	ResponseErrorAndCodeMessage(response, 500, err, err.Error())
}

func ResponseSystemError(response *restful.Response, err error) {
	ResponseErrorAndCodeMessage(response, 500, err, ErrMessageSystem)
}

//system error and error msg
func ResponseErrorAndMessage(response *restful.Response, err error, message string) {
	ResponseErrorAndCodeMessage(response, 500, err, message)
}
