package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager interface {
	NewToken() (string, error)
	Verify(token string) error
}

func NewJwtTokenManager(signingKey []byte, ttl time.Duration) TokenManager {
	return &jwtTokenManager{signingKey: signingKey, ttl: ttl}
}

type jwtTokenManager struct {
	signingKey []byte
	ttl        time.Duration
}

func (manager *jwtTokenManager) NewToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(manager.ttl).Unix(),
	})

	return token.SignedString(manager.signingKey)
}

func (manager *jwtTokenManager) Verify(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return manager.signingKey, nil
	})
	return err
}
