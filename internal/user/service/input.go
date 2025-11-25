package usersvc

import (
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

type UserCreateInput struct {
	UserBinding e.UserBinding
	Name        string
	DOB         time.Time
}

type UserUpdateInput struct {
	Name string
	DOB  time.Time
}
