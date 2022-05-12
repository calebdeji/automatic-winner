package handler

import "zate/generics"

type SignUpResponse struct {
	generics.GenericResponse
	Token string `json:"token"`
}

type VerifyEmailResponse struct {
	generics.GenericResponse
	Token         string      `json:"token"`
	Refresh_token string      `json:"refresh_token"`
	User_id       interface{} `json:"id"`
}
