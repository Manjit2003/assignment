package handler

import (
	"errors"
	"net/http"

	auth_service "github.com/Manjit2003/samespace/pkg/service/auth"
	"github.com/Manjit2003/samespace/pkg/utils"
)

//	@Summary		Login into your account
//	@Description	Returns the access token and refresh token upon successfull login
//	@Description	Please note that the accessToken will be valid only for 10 mins
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	utils.HTTPReponse
//	@Failure		500		{object}	utils.HTTPReponse
//	@Param			request	body		handler.HandleUserLogin.payload	true	"data"
//	@Router			/auth/login [post]
func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var data payload

	err := utils.ReadBody(w, r, &data)

	if err != nil {
		return
	}

	tokens, err := auth_service.LoginUser(data.Username, data.Password)

	if err != nil {
		if errors.Is(auth_service.ErrorInvalidCreds, err) {
			utils.SendResponse(w, 403, utils.HTTPReponse{
				Error:   true,
				Message: "invalid credentials",
			})
			return
		} else {
			utils.SendResponse(w, 500, utils.HTTPReponse{
				Error:   true,
				Message: "error loggin in user",
			})
			return
		}
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "login successfull",
		Data:    tokens,
	})
}

//	@Summary		Create new user account
//	@Description	Creates a new user account to add todos
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	utils.HTTPReponse
//	@Failure		500		{object}	utils.HTTPReponse
//	@Param			request	body		handler.HandleUserRegister.payload	true	"data"
//	@Router			/auth/register [post]
func HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var data payload

	err := utils.ReadBody(w, r, &data)

	if err != nil {
		return
	}

	err = auth_service.RegisterUser(data.Username, data.Password)

	if err != nil {
		if errors.Is(auth_service.ErrorUsernameExists, err) {
			utils.SendResponse(w, 403, utils.HTTPReponse{
				Error:   true,
				Message: "username already taken",
			})
			return
		} else {
			utils.SendResponse(w, 500, utils.HTTPReponse{
				Error:   true,
				Message: "error creating new account",
			})
			return
		}
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "register successfull",
	})
}

//	@Summary		Get new access token
//	@Description	Return's new access token from the refresh token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	utils.HTTPReponse
//	@Failure		500		{object}	utils.HTTPReponse
//	@Param			request	body		handler.HandleGetAccessToken.payload	true	"data"
//	@Router			/auth/refresh [post]
func HandleGetAccessToken(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	var data payload

	err := utils.ReadBody(w, r, &data)

	if err != nil {
		return
	}

	token, err := auth_service.GetAccessTokenFromRefreshToken(data.RefreshToken)

	if err != nil {
		utils.SendResponse(w, 500, utils.HTTPReponse{
			Error:   true,
			Message: "error getting access token",
		})
		return
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "token refreshed",
		Data: map[string]string{
			"access_token": token,
		},
	})
}
