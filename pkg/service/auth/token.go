package auth_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type JWTData struct {
	jwt.RegisteredClaims
	CustomClaims `json:"custom_claims"`
}

func generateAccessToken(userId string, username string) (string, error) {
	claims := JWTData{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		CustomClaims: CustomClaims{
			UserID:   userId,
			Username: username,
		},
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenString.SignedString([]byte("manjit_secret"))
	return token, err
}

func parseJWT(token string) (*CustomClaims, error) {
	claims := &JWTData{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("manjit_secret"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing the jwt: %v", err)
	}
	return &claims.CustomClaims, nil
}
