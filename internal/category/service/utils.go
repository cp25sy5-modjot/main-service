package category

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

func saveNewCategory(s *Service, cat *e.Category) (*m.CategoryRes, error) {
	newCat, err := s.repo.Create(cat)
	if err != nil {
		return nil, err
	}
	// Reload with preload
	catWithDetails, err := s.repo.FindByID(&m.CategorySearchParams{
		CategoryID: newCat.CategoryID,
		UserID:     newCat.UserID,
	})
	if err != nil {
		return nil, err
	}
	return buildCategoryResponse(catWithDetails), nil
}

func buildCategoryResponse(cat *e.Category) *m.CategoryRes {
	if cat.Transactions == nil {
		return &m.CategoryRes{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			Budget:       cat.Budget,
			ColorCode:    cat.ColorCode,
			CreatedAt:    cat.CreatedAt,
		}
	} else {
		budgetUsage := 0.0
		for _, tx := range cat.Transactions {
			_ = tx // just to avoid unused variable warning
			budgetUsage += tx.Price * tx.Quantity
		}
		return &m.CategoryRes{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			Budget:       cat.Budget,
			ColorCode:    cat.ColorCode,
			CreatedAt:    cat.CreatedAt,
			BudgetUsage:  budgetUsage,
		}
	}
}

func buildCategoryResponses(categories []e.Category) []m.CategoryRes {
	categoryResponses := make([]m.CategoryRes, 0, len(categories))
	for _, cat := range categories {
		res := buildCategoryResponse(&cat)
		categoryResponses = append(categoryResponses, *res)
	}
	return categoryResponses
}
