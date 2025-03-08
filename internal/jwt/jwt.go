/*
Package jwt implements a very simple JWT encoding method that satisfies GitHub's
authentication.

https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app
*/
package jwt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

const (
	iatDelta = -60
	expDelta = 600
)

var (
	ErrNotPEM     = errors.New("not PEM")
	ErrInvalidKey = errors.New("invalid key")
	ErrCannotSign = errors.New("cannot sign")
)

func New(now int64, id, key string) (string, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return "", ErrNotPEM
	}

	pKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidKey, err)
	}

	iat := now + iatDelta
	exp := now + expDelta

	header := encode(`{"typ":"JWT", "alg":"RS256"}`)
	payload := encode(
		fmt.Sprintf(
			`{"iat":%d,"exp":%d,"iss":"%s"}`,
			iat, exp, id,
		),
	)

	hash := sha256.New()
	hash.Write([]byte(header + "." + payload))
	sum := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(nil, pKey, crypto.SHA256, sum)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCannotSign, err)
	}

	token := header + "." + payload + "." + encode(string(signature))

	return token, nil
}

func encode(s string) string {
	return strings.TrimRight(
		base64.URLEncoding.EncodeToString([]byte(s)),
		"=",
	)
}
