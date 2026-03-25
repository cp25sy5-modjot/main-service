package model

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

type UserCreateInput struct {
	UserBinding e.UserBinding
	Name        string
}

type UserUpdateInput struct {
	Name string
}
