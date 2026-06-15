package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserId  uint  `json:"user_id"`
	Expired int64 `json:"expired"`
	jwt.RegisteredClaims
}

func GenToken(userId uint, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	expiredAt := now.Add(ttl)

	claims := TokenClaims{
		UserId:  userId,
		Expired: expiredAt.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Expired < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
