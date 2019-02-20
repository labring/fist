package main

import (
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	//	certutil "k8s.io/client-go/util/cert"
	certutil "k8s.io/client-go/util/cert"

	"fmt"
)

func createToken() {
	var alg jose.SignatureAlgorithm
	alg = jose.RS256

	privateKey := getPrivateKey(rsaPrivateKey)
	//	privateKey := getPublicKey(rsaPublicKey)

	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: alg,
			Key:       privateKey,
		},
		nil,
	)
	if err != nil {
		fmt.Errorf("wrong: %v", err)
	}

	type legacyPrivateClaimSa struct {
		ServiceAccountName string `json:"kubernetes.io/serviceaccount/service-account.name"`
		ServiceAccountUID  string `json:"kubernetes.io/serviceaccount/service-account.uid"`
		SecretName         string `json:"kubernetes.io/serviceaccount/secret.name"`
		Namespace          string `json:"kubernetes.io/serviceaccount/namespace"`
	}

	type legacyPrivateClaims struct {
		Aud  string
		Name string
	}

	//{
	//	"iss": "https://47.52.197.163:32000",
	//	"sub": "Cgc4OTEyNTU3EgZnaXRodWI",
	//	"aud": "example-app",
	//	"exp": 1524020707,
	//	"iat": 1523934307,
	//	"at_hash": "9s2hoIsPte1ogsuJzdZoZg",
	//	"email": "fhtjob@hotmail.com",
	//	"email_verified": true,
	//	"name": "steven"
	//}

	privateClaims := &legacyPrivateClaims{
		Aud:  "example-app",
		Name: "msxu",
	}

	// claims are applied in reverse precedence
	token, err := jwt.Signed(signer).
		Claims(privateClaims).
		Claims(&jwt.Claims{
			Issuer: "https://47.52.197.163:32000",
		}).
		CompactSerialize()
	if err != nil {
		fmt.Errorf("wrong: %v", err)
	}

	fmt.Println(token)

}

func getPrivateKey(data string) interface{} {
	key, _ := certutil.ParsePrivateKeyPEM([]byte(data))
	return key
}

const rsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAmsNkQ2sjfBbqg2p0FzAa875tEr++s2ARxrIbxbixoqlAeldd
os01eRM0hZz6nVvgOcjwK7IoQqgPGvdCIzF/iTzOkzL4K46e36Vfvy3M4fh9F1Cx
/Fn5/IoJWC2G4irOjlmZD+bxKNmTVya7OQkFjthOKUWsfTt6LIZpbFycNAAMQBLT
VLLkVbpkjb3JN+8iFJOjU/hqoy+yRbmDRSkpip0aBqVXniYqoveUOeIX71vHWY8M
lTtE/07U2qsIljDQTGXowtVONaUjnf/Np5QGy/EB6R8gWY7L0to+CLWhc8atC9uy
uEKu6BfCMYbcoqB+jnqg6Bp/MFtO0/PpUoGxMQIDAQABAoIBADoZXht9Lh4YkENz
hE9sLMISW+os949pYmMAXDK26mDRPzZuc+V5OjjQv+flDaRjaGLpD1ioEjsr0jfi
WP7TRzijDj3uZYIckYIOGEqyC/dNDNDi515//LwUqftjY/6l6VNBSZfRr/kQ7SJL
lP+NZnvAsl8GHAILgQsUDqGyhqVyRbIBMINonZcLk0+amelauvX3gr4xJwOFnDHy
MWCCiXkzcLa2PTNEifJNIP/7UO1OtxpHLj8ZHmEkxZIlxeOUNQNSmiR1J7hpmarE
NxVXtNpRiFgQTIP6VoCyv46e669AXehSjNbbLC+takqJ6rgZkdN3bz/TLVjCn22z
s4IzYsECgYEAyS9YCHhgd95AprEBXrmJK2L9QxeYobsBbdr4RSIKObePBV04+Yt8
ThsQi4nxfOP/0rfdJx9G1Tqr+e/fIXk61bNFExKbLJKzXsWU20e6pLKr84g1V/UM
8SmD1d30J3Ilm7U6RUEc2P8HKALpg/zDSNjybNTcqA2pWDEK5jv8XFkCgYEAxO4i
Kb1gHfPLqzQsDw2RUajPlkCkMnNddTuK0epZzKWlv5DzYFYCSGymf15EESXQIty1
4EarbbHaZrkGUUDWTJinU/SOHt+KiarAGtL8OnGwfmUm4FUJQ0vQCpwVqYPO/AAg
iUjoJHtd2d99tjDYBLM4KK3BPw2ionG5KqZSgJkCgYEAprxZHnP9p3qIZF4wf3wP
Rej9cxxcQDXHYm9m8YzboqgRnWV8cbz7oZPmq28At/wSKmZ9oA3Y26GLpFH7wDdD
3pZ7uenitxdCF1pzGyUgd50oy2Xop+QM/NXmUFpqHkMJDjotd/YV3XXHTY7UT7It
evNqP25PDex8m/3RRa0TYskCgYEAuKubK8sj5FKnzl95ZZBSkuIb8ImjsI/Qt0vj
zR/Xn/pCaVczc0aUk3gX1l7+s5njCZ4xjCSZQ5/B8AxYKUAE9gU8/JTb5YW5M4oE
5eKultvgJ1cR0tLLgekJKbne8nzhUB2KZVMSJovtoY9cIsrA9/9cjYELM+bEeVLt
0lnwChkCgYEAirjfoqUr9TCtWaJDIn10Dqxgk8wZn9w81Qn6VA3NAfei8zT4P/Lp
9ycF/YmCKjfsMqYlawahAgk2Of2fHUid5MjOPbBxo7i3DM4b9qi1GlmEkgmmiXOW
A0nS8mtlA3HR4GquUVgZSRbErgFmnMqSpGNyqRhgwt6pfK7Al7c/OTE=
-----END RSA PRIVATE KEY-----
`
