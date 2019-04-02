package auth

import (
	"github.com/emicklei/go-restful"
	"github.com/wonderivan/logger"
	"gopkg.in/square/go-jose.v2"
)

//key paires
var (
	Pub  jose.JSONWebKey
	Priv jose.JSONWebKey
)

//K8sRegister is k8s auth
func K8sRegister(auth *restful.WebService) {
	auth.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	//apiserver http
	auth.Route(auth.GET("/.well-known/openid-configuration").To(discoveryHandler))
	auth.Route(auth.GET("/keys").To(handlePublicKeys))
	//user token http
	auth.Route(auth.GET("/token").To(handlerToken))

}

func discoveryHandler(request *restful.Request, response *restful.Response) {
	type discovery struct {
		Issuer        string   `json:"issuer"`
		Auth          string   `json:"authorization_endpoint"`
		Token         string   `json:"token_endpoint"`
		Keys          string   `json:"jwks_uri"`
		ResponseTypes []string `json:"response_types_supported"`
		Subjects      []string `json:"subject_types_supported"`
		IDTokenAlgs   []string `json:"id_token_signing_alg_values_supported"`
		Scopes        []string `json:"scopes_supported"`
		AuthMethods   []string `json:"token_endpoint_auth_methods_supported"`
		Claims        []string `json:"claims_supported"`
	}

	dis := &discovery{
		Issuer:      "https://fist.sealyun.svc.cluster.local" + authHTTPSPortString,
		Auth:        "https://fist.sealyun.svc.cluster.local" + authHTTPSPortString + "/auth",
		Token:       "http://fist.sealyun.svc.cluster.local" + authHTTPPortString + "/token",
		Keys:        "https://fist.sealyun.svc.cluster.local" + authHTTPSPortString + "/keys",
		Subjects:    []string{"public"},
		IDTokenAlgs: []string{string(jose.RS256)},
		Scopes:      []string{"openid", "email", "groups", "profile", "offline_access"},
		AuthMethods: []string{"client_secret_basic"},
		Claims: []string{
			"aud", "email", "email_verified", "exp",
			"iat", "iss", "locale", "name", "sub", "groups",
		},
		ResponseTypes: []string{"code",
			"token",
			"id_token",
			"code token",
			"code id_token",
			"token id_token",
			"code token id_token",
			"none"},
	}

	logger.Info("discovery: %v", dis)
	_ = response.WriteEntity(dis)
}

func handlePublicKeys(request *restful.Request, response *restful.Response) {
	jwks := jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 1),
	}
	jwks.Keys[0] = Pub

	logger.Info("public keys: ", jwks)

	response.AddHeader("Content-Type", "application/json")
	_ = response.WriteEntity(&jwks)
}
