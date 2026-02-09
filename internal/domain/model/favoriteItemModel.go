package model

import "time"

type FavoriteItemInsertReq struct {
	Title      string  `json:"title" validate:"required,min=1,max=50"`
	Price      float64 `json:"price" validate:"required,gt=0"`
	CategoryID string  `json:"category_id" validate:"required,uuid4"`
}

type FavoriteItemUpdateReq struct {
	Title      *string  `json:"title" validate:"required,min=1,max=50"`
	Price      *float64 `json:"price" validate:"required,gt=0"`
	CategoryID *string  `json:"category_id" validate:"required,uuid4"`
}

type FavoriteItemReOrderReq struct {
	ReOrderList []FavoritePositionUpdateReq `json:"reorder_list" validate:"required,min=1,dive"`
}

type FavoritePositionUpdateReq struct {
	FavoriteID string `json:"favorite_id" validate:"required,uuid4"`
	Position   int    `json:"position" validate:"required,gte=0"`
}

type FavoriteItemRes struct {
	FavoriteID string    `json:"favorite_id"`
	Title      string    `json:"title"`
	Price      float64   `json:"price"`
	CategoryID string    `json:"category_id"`
	Position   int       `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FavoriteItemCreateInput struct {
	UserID     string
	Title      string
	Price      float64
	CategoryID string
}

type FavoriteItemUpdateInput struct {
	FavoriteID string
	UserID     string
	Title      *string
	Price      *float64
	CategoryID *string
}

type FavoritePositionUpdateInput struct {
	FavoriteID string
	Position   int
}

type FavoriteItemReOrderInput struct {
	UserID      string
	ReOrderList []FavoritePositionUpdateInput
}
