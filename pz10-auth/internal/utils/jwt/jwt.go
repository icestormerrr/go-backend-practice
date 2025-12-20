package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RS256TokenManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRS256TokenManager(privateKeyPEM, publicKeyPEM string) (*RS256TokenManager, error) {
	priv, err := parseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return nil, err
	}
	pub, err := parseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return nil, err
	}

	return &RS256TokenManager{
		privateKey: priv,
		publicKey:  pub,
	}, nil
}

func (r *RS256TokenManager) Sign(userID int64, email, role string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iat":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
		"iss":   "pz10-auth",
		"aud":   "pz10-clients",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(r.privateKey)
}

func (r *RS256TokenManager) Parse(tokenStr string) (map[string]any, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return r.publicKey, nil
	},
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithAudience("pz10-clients"),
		jwt.WithIssuer("pz10-auth"),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !t.Valid {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return map[string]any(claims), nil
}

func parseRSAPrivateKeyFromPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid RSA private key PEM")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func parseRSAPublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid RSA public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return pub, nil
}
