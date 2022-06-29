package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nooncall/owls/go/global"
)

func Login(userName, Pwd string) (string, error) {
	if err := loginService.Login(userName, Pwd); err != nil {
		return "", err
	}

	return GenerateToken(userName, Pwd)
}

type Claims struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	expireTime := time.Now().Add(time.Duration(global.GVA_CONFIG.Login.TokenEffectiveHour) * time.Hour)
	claims := Claims{username, password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ipalfish-db-injection",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenClaims.SignedString([]byte(global.GVA_CONFIG.Login.TokenSecret))
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.GVA_CONFIG.Login.TokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
		return nil, fmt.Errorf("parse token failed, not a claims ins")
	}
	return nil, fmt.Errorf("get nil token claims")
}
