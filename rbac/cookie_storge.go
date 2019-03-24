package rbac

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"
	"strings"
)

const desCookieKey = "df9gtsq3"

//generator token
func generatorCookieValue(userInfo *UserInfo) string {
	//serial
	userJSON, err := json.Marshal(userInfo)
	if err != nil {
		logger.Error("json format error:", err)
	}
	//des
	desVal := tools.DESEncrypt([]byte(string(userJSON)+"@@@@@@"+tools.MD5(string(userJSON))), []byte(desCookieKey))
	return desVal
}

// get user info from token
func getUserInfoFromToken(token string) *UserInfo {
	desVal := tools.DESDecrypt(token, []byte(desCookieKey))
	data := strings.Split(desVal, "@@@@@@")
	if len(data) == 2 {
		md5Val := data[1]
		serialVal := data[0]
		if tools.MD5(serialVal) == md5Val {
			userInfo := &UserInfo{}
			err := json.Unmarshal([]byte(serialVal), userInfo)
			if err != nil {
				logger.Error(err)
			} else {
				return userInfo
			}
		} else {
			logger.Error("the md5 value not equals serial data")
		}
	} else {
		logger.Error("validateCookieToken data spilt length not equals 2")
	}
	return nil
}

// validate token and username
func validateCookieToken(token, username string) bool {
	userInfo := getUserInfoFromToken(token)
	if userInfo != nil {
		return username == userInfo.Username
	}
	return false
}

// used login restful
func loginCookieSetter(response *restful.Response, userInfo *UserInfo) {
	//1. username
	tools.CookieWriteValue(response, "username", userInfo.Username)
	//2. login status
	tools.CookieWriteValue(response, "logged", "yes")
	//3. token
	tools.CookieWriteValue(response, "user_token", generatorCookieValue(userInfo))
}

// used logout restful
func logoutCookieSetter(response *restful.Response) {
	//1. username
	tools.CookieRemoveValue(response, "username")
	//2. login status
	tools.CookieRemoveValue(response, "logged")
	//3. token
	tools.CookieRemoveValue(response, "user_token")
}

// used filter validate
func filterCookieValidate(req *restful.Request) bool {
	logged := tools.CookieRead(req, "logged") // is 'yes'
	token := tools.CookieRead(req, "user_token")
	username := tools.CookieRead(req, "username")
	if "yes" == logged && validateCookieToken(token, username) {
		return true
	}
	return false
}

//FistCookieGetUserInfo is get fist cookie user info
func FistCookieGetUserInfo(req *restful.Request) *UserInfo {
	if filterCookieValidate(req) {
		//if validated
		token := tools.CookieRead(req, "user_token")
		return getUserInfoFromToken(token)
	}
	return nil
}

//FistCookieUpdateUserInfo is update fist cookie user info
func FistCookieUpdateUserInfo(response *restful.Response, userInfo *UserInfo) {
	loginCookieSetter(response, userInfo)
}
