package usermodel

type AuthenticateRes struct {
	Token string `json:"token"`
	ExpIn int    `json:"expIn"`
}
