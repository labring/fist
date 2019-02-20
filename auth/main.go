package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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

// Determine the signature algorithm for a JWT.
func signatureAlgorithm(jwk *jose.JSONWebKey) (alg jose.SignatureAlgorithm, err error) {
	if jwk.Key == nil {
		return alg, errors.New("no signing key")
	}
	switch key := jwk.Key.(type) {
	case *rsa.PrivateKey:
		// Because OIDC mandates that we support RS256, we always return that
		// value. In the future, we might want to make this configurable on a
		// per client basis. For example allowing PS256 or ECDSA variants.
		//
		// See https://github.com/coreos/dex/issues/692
		return jose.RS256, nil
	case *ecdsa.PrivateKey:
		// We don't actually support ECDSA keys yet, but they're tested for
		// in case we want to in the future.
		//
		// These values are prescribed depending on the ECDSA key type. We
		// can't return different values.
		switch key.Params() {
		case elliptic.P256().Params():
			return jose.ES256, nil
		case elliptic.P384().Params():
			return jose.ES384, nil
		case elliptic.P521().Params():
			return jose.ES512, nil
		default:
			return alg, errors.New("unsupported ecdsa curve")
		}
	default:
		return alg, fmt.Errorf("unsupported signing key type %T", key)
	}
}

type idTokenClaims struct {
	Issuer           string `json:"iss"`
	Subject          string `json:"sub"`
	Audience         string `json:"aud"`
	Expiry           int64  `json:"exp"`
	IssuedAt         int64  `json:"iat"`
	AuthorizingParty string `json:"azp,omitempty"`
	Nonce            string `json:"nonce,omitempty"`

	AccessTokenHash string `json:"at_hash,omitempty"`

	Email         string `json:"email,omitempty"`
	EmailVerified *bool  `json:"email_verified,omitempty"`

	Groups []string `json:"groups,omitempty"`

	Name string `json:"name,omitempty"`

	FederatedIDClaims *federatedIDClaims `json:"federated_claims,omitempty"`
}

type federatedIDClaims struct {
	ConnectorID string `json:"connector_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

func signPayload(key *jose.JSONWebKey, alg jose.SignatureAlgorithm, payload []byte) (jws string, err error) {
	signingKey := jose.SigningKey{Key: key, Algorithm: alg}

	signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		return "", fmt.Errorf("new signier: %v", err)
	}
	signature, err := signer.Sign(payload)
	if err != nil {
		return "", fmt.Errorf("signing payload: %v", err)
	}
	return signature.CompactSerialize()
}

func handlerToken(request *restful.Request, response *restful.Response) {
	signingAlg, err := signatureAlgorithm(&Priv)
	if err != nil {
		return
	}

	ev := true

	tok := idTokenClaims{
		Issuer:        "https://dex.example.com:8080",
		Subject:       "Cgc4OTEyNTU3EgZnaXRodWI",
		Audience:      "example-app",
		Expiry:        time.Now().Add(time.Minute * 10).Unix(),
		IssuedAt:      time.Now().Unix(),
		Email:         "fhtjob@hotmail.com",
		EmailVerified: &ev,
		Groups:        []string{"dev"},
		Name:          "fanux",
	}

	payload, err := json.Marshal(tok)
	if err != nil {
		fmt.Println("could not serialize claims", err)
		return
	}

	var idToken string
	if idToken, err = signPayload(&Priv, signingAlg, payload); err != nil {
		fmt.Println("failed to sign payload", err)
		return
	}

	response.WriteEntity(&idToken)
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
