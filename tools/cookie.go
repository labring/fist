package tools

import (
	"github.com/emicklei/go-restful"
	"github.com/wonderivan/logger"
	"net/http"
)

//CookieWriteValue is write cookies  of name value
func CookieWriteValue(resp *restful.Response, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}
	http.SetCookie(resp, &cookie)
}

//CookieRemoveValue is remove cookies  of name value
func CookieRemoveValue(resp *restful.Response, name string) {
	cookie := http.Cookie{
		Name:   name,
		MaxAge: int(-1),
	}
	http.SetCookie(resp, &cookie)
}

//CookieRead is read value from cookie
func CookieRead(req *restful.Request, name string) string {
	cookie, err := req.Request.Cookie(name)
	if err != nil {
		logger.Error("cookie read is error,", err, "; cookie name is :", name)
		return ""
	}
	return cookie.Value
}
