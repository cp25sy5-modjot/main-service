package mapper

import (
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

func ParseFavoriteItemInsertReqToServiceInput(
	userID string,
	req *m.FavoriteItemInsertReq,
) *m.FavoriteItemCreateInput {
	return &m.FavoriteItemCreateInput{
		UserID:     userID,
		Title:      req.Title,
		CategoryID: req.CategoryID,
		Price:      req.Price,
	}
}

func BuildFavoriteItemResponse(fav *e.FavoriteItem) *m.FavoriteItemRes {
	return &m.FavoriteItemRes{
		FavoriteID: fav.FavoriteID,
		Title:      fav.Title,
		CategoryID: fav.CategoryID,
		Price:      fav.Price,
		Position:   fav.Position,
		CreatedAt:  fav.CreatedAt,
		UpdatedAt:  fav.UpdatedAt,

		CategoryIcon:  fav.Category.Icon,
		CategoryColor: fav.Category.ColorCode,
	}
}

func BuildFavObjectToCreate(fav_id string, input *m.FavoriteItemCreateInput) *e.FavoriteItem {
	return &e.FavoriteItem{
		FavoriteID: fav_id,
		UserID:     input.UserID,
		Title:      input.Title,
		CategoryID: input.CategoryID,
		Price:      input.Price,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ParseFavoriteItemUpdateReqToServiceInput(
	userID string,
	favID string,
	req *m.FavoriteItemUpdateReq,
) *m.FavoriteItemUpdateInput {
	return &m.FavoriteItemUpdateInput{
		UserID:     userID,
		FavoriteID: favID,
		Title:      req.Title,
		CategoryID: req.CategoryID,
		Price:      req.Price,
	}
}

func ParseFavoriteItemReOrderReqToServiceInput(
	userID string,
	req *m.FavoriteItemReOrderReq,
) *m.FavoriteItemReOrderInput {
	var reorderList []m.FavoritePositionUpdateInput
	for _, item := range req.ReOrderList {
		reorderList = append(reorderList, m.FavoritePositionUpdateInput{
			FavoriteID: item.FavoriteID,
			Position:   item.Position,
		})
	}

	return &m.FavoriteItemReOrderInput{
		UserID:      userID,
		ReOrderList: reorderList,
	}
}
