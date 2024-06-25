package model

type AuthLoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRegisterData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRegisterResponse struct {
	Success bool `json:"success"`
}
