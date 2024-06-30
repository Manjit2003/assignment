package auth_service_test

import (
	"testing"

	"github.com/Manjit2003/samespace/pkg/config"
	"github.com/Manjit2003/samespace/pkg/db"
	auth_service "github.com/Manjit2003/samespace/pkg/service/auth"
	"github.com/Manjit2003/samespace/pkg/utils"
)

func TestUserRegister(t *testing.T) {
	db.InitDatabase(&config.TestConfig)

	user, pass := utils.GenerateRandomCreds()

	err := auth_service.RegisterUser(user, pass)

	if err != nil {
		t.Errorf("error registering new user: %v", err)
	}
}

func TestLoginUser(t *testing.T) {
	db.InitDatabase(&config.TestConfig)

	user, pass := utils.GenerateRandomCreds()

	err := auth_service.RegisterUser(user, pass)

	if err != nil {
		t.Errorf("error registering new user: %v", err)
	}

	tokens, err := auth_service.LoginUser(user, pass)

	if err != nil {
		t.Errorf("error in user login: %v", err)
	}

	if tokens.JWTToken == "" {
		t.Errorf("returned blank jwt")
	}

	if tokens.RefreshToken == "" {
		t.Errorf("returned blank refresh token")
	}
}

func TestRefreshToken(t *testing.T) {
	db.InitDatabase(&config.TestConfig)

	user, pass := utils.GenerateRandomCreds()

	err := auth_service.RegisterUser(user, pass)

	if err != nil {
		t.Errorf("error registering new user: %v", err)
	}

	tokens, err := auth_service.LoginUser(user, pass)

	if err != nil {
		t.Errorf("error in user login: %v", err)
	}

	if tokens.JWTToken == "" {
		t.Errorf("returned blank jwt")
	}

	if tokens.RefreshToken == "" {
		t.Errorf("returned blank refresh token")
	}

	newToken, err := auth_service.GetAccessTokenFromRefreshToken(string(tokens.RefreshToken))
	if err != nil {
		t.Errorf("error getting jwt from refresh token: %v", err)
	}

	if newToken == "" {
		t.Errorf("new token is blank")
	}
}
