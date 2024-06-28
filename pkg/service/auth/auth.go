package auth_service

import (
	"fmt"
	"time"

	user_service "github.com/Manjit2003/samespace/pkg/service/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorInvalidCreds   = fmt.Errorf("invalid credentials")
	ErrorUsernameExists = fmt.Errorf("username already exists")
)

type JWTToken string
type RefreshToken string

type AuthTokens struct {
	JWTToken
	RefreshToken
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginUser(username, password string) (*AuthTokens, error) {
	user, err := user_service.GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	if user == nil {
		return nil, ErrorInvalidCreds
	}
	if !checkPasswordHash(password, user.HashedPassword) {
		return nil, ErrorInvalidCreds
	}
	authToken, err := generateToken(user.ID, user.Username, time.Minute*10)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %v", err)
	}
	refreshToken, err := generateToken(user.ID, user.Username, time.Hour*24*7) // a week
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %v", err)
	}
	return &AuthTokens{
		JWTToken:     JWTToken(authToken),
		RefreshToken: RefreshToken(refreshToken),
	}, nil
}

func RegisterUser(username, password string) error {
	user, err := user_service.GetUser(username)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}
	if user != nil {
		return ErrorUsernameExists
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing the user password: %v", err)
	}
	return user_service.AddUser(username, hashed)
}

func GetAccessTokenFromRefreshToken(token string) (string, error) {
	claims, err := parseJWT(token)
	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}
	return generateToken(claims.UserID, claims.Username, time.Minute*10)
}

func GetUserFromJWT(token string) (string, error) {
	claims, err := parseJWT(token)
	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}
	return claims.UserID, nil
}
