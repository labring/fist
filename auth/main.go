package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	jose "gopkg.in/square/go-jose.v2"
)

// This example has the same service definition as restful-user-resource
// but uses a different router (CurlyRouter) that does not use regular expressions
//
// POST http://localhost:8080/users
// <User><ID>1</ID><Name>Melissa Raspberry</Name></User>
//
// GET http://localhost:8080/users/1
//
// PUT http://localhost:8080/users/1
// <User><ID>1</ID><Name>Melissa</Name></User>
//
// DELETE http://localhost:8080/users/1
//

//User is
type User struct {
	ID, Name string
}

//key paires
var (
	Pub  jose.JSONWebKey
	Priv jose.JSONWebKey
)

//UserResource is
type UserResource struct {
	// normally one would use DAO (data access object)
	users map[string]User
}

//Register is
func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/{user-id}").To(u.findUser))
	ws.Route(ws.POST("").To(u.updateUser))
	ws.Route(ws.PUT("/{user-id}").To(u.createUser))
	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser))

	container.Add(ws)

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
		Issuer:      "https://dex.example.com:8080",
		Auth:        "https://dex.example.com:8080/auth",
		Token:       "https://dex.example.com:8080/token",
		Keys:        "https://dex.example.com:8080/keys",
		Subjects:    []string{"public"},
		IDTokenAlgs: []string{string(jose.RS256)},
		Scopes:      []string{"openid", "email", "groups", "profile", "offline_access"},
		AuthMethods: []string{"client_secret_basic"},
		Claims: []string{
			"aud", "email", "email_verified", "exp",
			"iat", "iss", "locale", "name", "sub",
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

	response.WriteEntity(dis)
}

func handlerToken(request *restful.Request, response *restful.Response) {
}

func newUUID() string {
	u := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, u); err != nil {
		panic(err)
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u)
}

//CreateKeyPair is
func CreateKeyPair() (priv, pub jose.JSONWebKey) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("gen rsa key: %v", err)
	}
	priv = jose.JSONWebKey{
		Key:       key,
		KeyID:     newUUID(),
		Algorithm: "RS256",
		Use:       "sig",
	}
	pub = jose.JSONWebKey{
		Key:       key.Public(),
		KeyID:     newUUID(),
		Algorithm: "RS256",
		Use:       "sig",
	}

	return priv, pub
}

func handlePublicKeys(request *restful.Request, response *restful.Response) {
	jwks := jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 1),
	}
	jwks.Keys[0] = Pub
	//TODO VerificationKeys

	response.WriteEntity(&jwks)
}

// GET http://localhost:8080/users/1
//
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	response.WriteEntity(&User{id, "fanux"})
}

// PUT http://localhost:8080/users/1
// <User><ID>1</ID><Name>Melissa</Name></User>
//
func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = *usr
		response.WriteEntity(usr)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// POST http://localhost:8080/users
// <User><ID>1</ID><Name>Melissa Raspberry</Name></User>
//
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := User{ID: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = usr
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// DELETE http://localhost:8080/users/1
//
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	u := UserResource{map[string]User{}}
	u.Register(wsContainer)

	log.Print("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServeTLS("ssl/cert.pem", "ssl/key.pem"))
}
