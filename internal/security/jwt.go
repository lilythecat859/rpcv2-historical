// SPDX-License-Identifier: AGPL-3.0-only
package security

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Signer interface {
	Sign(method string) (string, error)
}

type Validator interface {
	Validate(bearer string) (scope uint32, err error)
}

type edKey struct {
	mu   sync.RWMutex
	priv ed25519.PrivateKey
	pub  ed25519.PublicKey
}

func NewEdKey() (*edKey, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &edKey{priv: priv, pub: pub}, nil
}

func ParseEdKey(b64 string) (*edKey, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(b) != ed25519.PrivateKeySize {
		return nil, errors.New("bad private key size")
	}
	k := &edKey{priv: ed25519.PrivateKey(b)}
	k.pub = k.priv.Public().(ed25519.PublicKey)
	return k, nil
}

func (k *edKey) Sign(method string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   method,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(k.priv)
}

func (k *edKey) Validate(bearer string) (uint32, error) {
	k.mu.RLock()
	pub := k.pub
	k.mu.RUnlock()
	token, err := jwt.Parse(bearer, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEdDSA); !ok {
			return nil, errors.New("bad algo")
		}
		return pub, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	return parseScope(sub), nil
}

func (k *edKey) Rotate() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	k.priv = priv
	k.pub = pub
	return nil
}

func parseScope(sub string) uint32 {
	switch sub {
	case "read":
		return ScopeRead
	case "write":
		return ScopeWrite | ScopeRead
	case "admin":
		return ScopeRead | ScopeWrite | ScopeAdmin
	}
	return 0
}
