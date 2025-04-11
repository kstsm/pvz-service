package models

import "github.com/google/uuid"

type UserRegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserRegisterResp struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResp struct {
	Password string `json:"password"`
	Role     string `json:"role"`
}

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type DummyLoginResponse struct {
	Token string `json:"token"`
}
