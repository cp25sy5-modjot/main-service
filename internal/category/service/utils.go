package category

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	r "github.com/cp25sy5-modjot/main-service/internal/response/error"
)

func validateCategoryOwnership(cate *e.Category, userID string) error {
	if cate.UserID != userID {
		return r.Conflict(nil, "You are not authorized to access this category")
	}
	return nil
}
