package user

import (
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
)

func buildUserObjectToCreate(uid string, user *m.UserInsertReq) *e.User {
	return &e.User{
		UserID: uid,
		UserBinding: e.UserBinding{
			GoogleID:   user.UserBinding.GoogleID,
			FacebookID: user.UserBinding.FacebookID,
			AppleID:    user.UserBinding.AppleID,
		},
		Name:      user.Name,
		DOB:       user.DOB,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createDefaultCategories(s *Service, uid string) error {
	defaultCategories := []string{"อาหาร", "การเดินทาง", "ความบันเทิง", "ชอปปิ้ง", "อื่นๆ"}
	for _, categoryName := range defaultCategories {
		_, err := s.cat.Create(&e.Category{
			CategoryName: categoryName,
			UserID:       uid,
			Budget:       1000.0,
			ColorCode:    utils.GenerateRandomColor(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
