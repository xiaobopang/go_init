package lib

import (
	jwt "github.com/dgrijalva/jwt-go"
)

var config ServerConfig

func init() {
	config = LoadServerConfig()
}

var jwtSecret = []byte(config.JwtSecret)

type Claims struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}

func GenerateToken(id int, email string) (string, error) {
	claims := Claims{
		id,
		email,
		jwt.StandardClaims{
			ExpiresAt: config.JwtTokenExpire,
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
