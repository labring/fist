package tools

import (
	"github.com/emicklei/go-restful"
	"github.com/wonderivan/logger"
)

type responseObject struct {
	Message string      `json:"message"`
	Code    int32       `json:"code"`
	Data    interface{} `json:"data"`
}

//ResponseSuccess is web response for success
func ResponseSuccess(response *restful.Response, data interface{}) {
	logger.Info("return response data: ", data)
	err := response.WriteEntity(responseObject{Code: 200, Message: "success", Data: data})
	if err != nil {
		logger.Error("return response error: ", err)
	}
}

//ResponseErrorAndCodeMessage is web response for error
func ResponseErrorAndCodeMessage(response *restful.Response, code int32, err error, message string) {
	logger.Error("response error: ", err)
	err = response.WriteEntity(responseObject{Code: code, Message: message, Data: ""})
	if err != nil {
		logger.Error("return response error: ", err)
	}
}

//ResponseError is web response for error
func ResponseError(response *restful.Response, err error) {
	ResponseErrorAndCodeMessage(response, 500, err, err.Error())
}

//ResponseSystemError is web response for error
func ResponseSystemError(response *restful.Response, err error) {
	ResponseErrorAndCodeMessage(response, 500, err, ErrMessageSystem)
}

//ResponseErrorAndMessage is web response for error
func ResponseErrorAndMessage(response *restful.Response, err error, message string) {
	ResponseErrorAndCodeMessage(response, 500, err, message)
}

//ResponseAuthError is web response for error of cookie not login
func ResponseAuthError(response *restful.Response) {
	ResponseErrorAndCodeMessage(response, 401, ErrUserNotLogin, ErrUserNotLogin.Error())
}
