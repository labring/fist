package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/fanux/fist/tools"
	"gopkg.in/square/go-jose.v2"
)

//loadKeyPair is
func loadKeyPair() (pub, priv jose.JSONWebKey, err error) {

	privKey, _ := tools.ParseRsaPrivateKeyFromPemFile(PrivateKey)
	publicKey, _ := tools.ParseRsaPubKeyFromPemFile(PublicKey)
	//if file is not pem file format
	if privKey == nil || publicKey == nil {
		privKey = tools.PemDefaultPrivateKey()
		publicKey = tools.PemDefaultPublicKey()
		if privKey == nil || publicKey == nil {
			return jose.JSONWebKey{}, jose.JSONWebKey{}, tools.ErrAuthTokeKeyError
		}
	}
	priv = jose.JSONWebKey{
		Key:       privKey,
		KeyID:     "Cgc4OTEyNTU3EgZnaXRodWI",
		Algorithm: "RS256",
		Use:       "sig",
	}
	pub = jose.JSONWebKey{
		Key:       publicKey,
		KeyID:     "Cgc4OTEyNTU3EgZnaXRodWI",
		Algorithm: "RS256",
		Use:       "sig",
	}

	return pub, priv, nil
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
