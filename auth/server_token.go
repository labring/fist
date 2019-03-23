package auth

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/fanux/fist/tools"
	"github.com/wonderivan/logger"
	"time"
)

//TokenRegister is k8s auth token
func TokenRegister(auth *restful.WebService) {
	Pub, Priv = CreateKeyPair()
	auth.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	//user token http
	auth.Route(auth.GET("/token").To(handlerToken))

}

func handlerToken(request *restful.Request, response *restful.Response) {
	groups := request.Request.URL.Query()["group"]
	user := request.QueryParameter("user")
	logger.Info("user: ", user, ", groups: ", groups, ", url value:", request.Request.URL.Query())

	signingAlg, err := signatureAlgorithm(&Priv)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}

	ev := true
	tok := idTokenClaims{
		Issuer:        "https://fist.sealyun.svc.cluster.local" + authHTTPSPortString,
		Subject:       "Cgc4OTEyNTU3EgZnaXRodWI",
		Audience:      "sealyun-fist",
		Expiry:        time.Now().Add(time.Hour * 100).Unix(),
		IssuedAt:      time.Now().Unix(),
		Email:         "fhtjob@hotmail.com",
		EmailVerified: &ev,
		Groups:        groups,
		Name:          user,
	}

	payload, err := json.Marshal(&tok)
	logger.Info("token claims: %s", payload)
	if err != nil {
		tools.ResponseSystemError(response, err)
		return
	}

	var idToken string
	if idToken, err = signPayload(&Priv, signingAlg, payload); err != nil {
		tools.ResponseSystemError(response, err)
		return
	}

	logger.Info("token: ", idToken)
	tools.ResponseSuccess(response, &idToken)
}
