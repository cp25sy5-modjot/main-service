package model

type PushRegisterReq struct {
	Token    string `json:"token" validate:"required"`
	Platform string `json:"platform"`
}