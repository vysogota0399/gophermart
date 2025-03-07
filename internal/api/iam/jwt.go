package iam

import (
	"fmt"

	"github.com/vysogota0399/gophermart_portal/internal/api/models"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Sid uint64
}

func (i *Iam) decode(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(i.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("internal/api/iam/jwt: decode token failed error %w", err)
	}

	return claims, nil
}

func (i *Iam) buildJWTString(session *models.Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(session.ExpiredAt),
			Subject:   session.Sub,
			IssuedAt:  jwt.NewNumericDate(session.CreatedAt),
		},
		Sid: session.ID,
	})

	tokenString, err := token.SignedString([]byte(i.SecretKey))
	if err != nil {
		return "", fmt.Errorf("internal/api/iam/jwt: sign token failed error %w", err)
	}

	return tokenString, nil
}
