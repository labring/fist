package rbac

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
	//GET_USER ALL
	auth.Route(auth.GET("/user").To(handleListUserInfo))
	//GET_USER SINGLE
	auth.Route(auth.GET("/user/{user_name}").To(handleGetUserInfo))
	//ADD_USER
	auth.Route(auth.POST("/user").To(handleAddUserInfo))
	//UPDATE_USER
	auth.Route(auth.PUT("/user").To(handleUpdateUserInfo))
	//DELETE_USER
	auth.Route(auth.DELETE("/user/{user_name}").To(handleDelUserInfo))
}

func handleLogin(request *restful.Request, response *restful.Response) {
	t := &UserInfo{}
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	uerInfo := DoAuthentication(t.Username, t.Password)
	if uerInfo == nil {
		tools.ResponseError(response, tools.ErrUserAuth)
		return
	}
	tools.ResponseSuccess(response, uerInfo)
}

func handleGetUserInfo(request *restful.Request, response *restful.Response) {
	userName := request.PathParameter("user_name")
	userInfo := GetUserInfo(userName)
	if userInfo == nil {
		tools.ResponseError(response, tools.ErrUserGet)
		return
	}
	tools.ResponseSuccess(response, userInfo)
}

func handleListUserInfo(request *restful.Request, response *restful.Response) {
	arr := ListAllUserInfo()
	tools.ResponseSuccess(response, arr)
}

func handleAddUserInfo(request *restful.Request, response *restful.Response) {
	t := &UserInfo{}
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	err = AddUserInfo(t)
	if err != nil {
		tools.ResponseError(response, tools.ErrUserAdd)
		return
	}
	tools.ResponseSuccess(response, nil)
}

func handleUpdateUserInfo(request *restful.Request, response *restful.Response) {
	t := &UserInfo{}
	err := request.ReadEntity(t)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}
	err = UpdateUserInfo(t)
	if err != nil {
		tools.ResponseError(response, tools.ErrUserAdd)
		return
	}
	tools.ResponseSuccess(response, nil)
}

func handleDelUserInfo(request *restful.Request, response *restful.Response) {
	userName := request.PathParameter("user_name")
	err := DelUserInfo(userName)
	if err != nil {
		tools.ResponseError(response, tools.ErrUserDel)
		return
	}
	tools.ResponseSuccess(response, nil)
}
