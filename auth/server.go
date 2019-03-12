package auth

import (
	"encoding/json"
	"fmt"
	"github.com/fanux/fist/tools"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	jose "gopkg.in/square/go-jose.v2"
)

//key paires
var (
	Pub  jose.JSONWebKey
	Priv jose.JSONWebKey
)

//Register is
func Register(container *restful.Container) {
	Pub, Priv = CreateKeyPair()

	auth := new(restful.WebService)
	auth.
		Path("/").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	auth.Route(auth.GET("/.well-known/openid-configuration").To(discoveryHandler))
	auth.Route(auth.GET("/token").To(handlerToken))
	auth.Route(auth.GET("/keys").To(handlePublicKeys))

	container.Add(auth)

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
		Issuer:      "https://fist.sealyun.svc.cluster.local:8080",
		Auth:        "https://fist.sealyun.svc.cluster.local:8080/auth",
		Token:       "https://fist.sealyun.svc.cluster.local:8080/token",
		Keys:        "https://fist.sealyun.svc.cluster.local:8080/keys",
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

	fmt.Printf("discovery: %v", dis)
	response.WriteEntity(dis)
}

func handlerToken(request *restful.Request, response *restful.Response) {
	groups := request.Request.URL.Query()["group"]
	user := request.QueryParameter("user")
	fmt.Printf("user: %s, groups: %v, url value: %v", user, groups, request.Request.URL.Query())

	signingAlg, err := signatureAlgorithm(&Priv)
	if err != nil {
		fmt.Println("failed to sign payload", err)
		tools.ResponseError(response, tools.ErrSignPayload)
		return
	}

	ev := true
	tok := idTokenClaims{
		Issuer:        "https://fist.sealyun.svc.cluster.local:8080",
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
	fmt.Printf("token claims: %s", payload)
	if err != nil {
		fmt.Println("could not serialize claims", err)
		tools.ResponseError(response, tools.ErrSerializeClaims)
		return
	}

	var idToken string
	if idToken, err = signPayload(&Priv, signingAlg, payload); err != nil {
		fmt.Println("failed to sign payload", err)
		tools.ResponseError(response, tools.ErrSignPayload)
		return
	}

	fmt.Println("token: ", idToken)
	tools.ResponseSuccess(response, &idToken)
}

func handlePublicKeys(request *restful.Request, response *restful.Response) {
	jwks := jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 1),
	}
	jwks.Keys[0] = Pub

	fmt.Printf("public keys: %v", jwks)

	response.AddHeader("Content-Type", "application/json")
	response.WriteEntity(&jwks)
}

//Serve start a auth server
func Serve() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	Register(wsContainer)

	log.Print("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServeTLS("/etc/fist/cert.pem", "/etc/fist/key.pem"))
}
