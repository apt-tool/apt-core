package jwt

import (
	"time"

	"github.com/automated-pen-testing/api/pkg/enum"

	"github.com/golang-jwt/jwt/v4"
)

type authenticator struct {
	key    string
	expire int
}

func (a *authenticator) GenerateToken(name string, role enum.Role) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expireTime := time.Now().Add(time.Duration(a.expire) * time.Minute)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = expireTime.Unix()
	claims["name"] = name
	claims["role"] = role

	tokenString, err := token.SignedString([]byte(a.key))
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expireTime, nil
}

func (a *authenticator) ParseToken(token string) (string, error) {
	tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", ErrSigningMethod
		}

		return []byte(a.key), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := tokenString.Claims.(jwt.MapClaims); ok && tokenString.Valid {
		name := claims["name"].(string)

		return name, nil
	}

	return "", ErrInvalidToken
}
