package auth

import (
	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
)

//FistRegister is fist auth controller
func FistRegister(auth *restful.WebService) {
	auth.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	//login http server
	auth.Route(auth.POST("/login").To(handleLogin))
	//user manager
	//TODO not finish
	//GET_USER ALL
	auth.Route(auth.GET("/user").To(handleLogin))
	//GET_USER SINGLE
	auth.Route(auth.GET("/user/{user_name}").To(handleLogin))
	//ADD_USER
	auth.Route(auth.POST("/user").To(handleLogin))
	//UPDATE_USER
	auth.Route(auth.PUT("/user").To(handleLogin))
	//DELETE_USER
	auth.Route(auth.DELETE("/user/{user_name}").To(handleLogin))
}

func handleLogin(request *restful.Request, response *restful.Response) {
	t := &UserInfo{}
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	uerInfo := DoAuthentication(t.Name, t.Password)
	if uerInfo == nil {
		tools.ResponseError(response, tools.ErrUserAuth)
		return
	}
	tools.ResponseSuccess(response, uerInfo)
}
