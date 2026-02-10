package favrepo

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(favorite *e.FavoriteItem) (*e.FavoriteItem, error) {
	if err := r.db.Create(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *Repository) CreateTx(
	tx *gorm.DB,
	favorite *e.FavoriteItem,
) (*e.FavoriteItem, error) {
	if err := tx.Create(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *Repository) FindAll(uid string) ([]*e.FavoriteItem, error) {
	var favorites []*e.FavoriteItem
	err := r.db.
		Where("user_id = ?", uid).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("category_id", "icon", "color_code")
		}).
		Order("position ASC").
		Find(&favorites).Error
	return favorites, err
}

func (r *Repository) FindByID(uid, favoriteID string) (*e.FavoriteItem, error) {
	var favorite e.FavoriteItem
	err := r.db.
		Where("favorite_id = ? AND user_id = ?", favoriteID, uid).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("category_id", "icon", "color_code")
		}).
		First(&favorite).Error
	return &favorite, err
}

func (r *Repository) FindByIDTx(
	tx *gorm.DB,
	uid,
	favoriteID string,
) (*e.FavoriteItem, error) {
	var favorite e.FavoriteItem
	err := tx.
		Where("favorite_id = ? AND user_id = ?", favoriteID, uid).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("category_id", "icon", "color_code")
		}).
		First(&favorite).Error
	return &favorite, err
}

func (r *Repository) Update(favorite *e.FavoriteItem) (*e.FavoriteItem, error) {
	if err := r.db.Save(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *Repository) Delete(uid, favoriteID string) error {
	return r.db.
		Delete(&e.FavoriteItem{}, "favorite_id = ? AND user_id = ?", favoriteID, uid).
		Error
}

func (r *Repository) DeleteTx(
	tx *gorm.DB,
	uid,
	favoriteID string,
) error {
	return tx.
		Delete(&e.FavoriteItem{}, "favorite_id = ? AND user_id = ?", favoriteID, uid).
		Error
}

func (r *Repository) GetMaxPosition(uid string) (int, error) {
	var max int
	err := r.db.Model(&e.FavoriteItem{}).
		Where("user_id = ?", uid).
		Select("COALESCE(MAX(position), 0)").
		Scan(&max).Error
	return max, err
}

func (r *Repository) GetMaxPositionTx(
	tx *gorm.DB,
	uid string,
) (int, error) {
	var max int
	err := tx.Model(&e.FavoriteItem{}).
		Where("user_id = ?", uid).
		Select("COALESCE(MAX(position), 0)").
		Scan(&max).Error
	return max, err
}

func (r *Repository) ShiftLeftAfter(uid string, pos int) error {
	return r.db.Model(&e.FavoriteItem{}).
		Where("user_id = ? AND position > ?", uid, pos).
		Update("position", gorm.Expr("position - 1")).
		Error
}

func (r *Repository) ShiftLeftAfterTx(
	tx *gorm.DB,
	uid string,
	pos int,
) error {
	return tx.Model(&e.FavoriteItem{}).
		Where("user_id = ? AND position > ?", uid, pos).
		Update("position", gorm.Expr("position - 1")).
		Error
}

func (r *Repository) UpdatePosition(
	uid string,
	favoriteID string,
	position int,
) error {
	return r.db.Model(&e.FavoriteItem{}).
		Where("favorite_id = ? AND user_id = ?", favoriteID, uid).
		Update("position", position).
		Error
}

func (r *Repository) UpdatePositionTx(
	tx *gorm.DB,
	uid string,
	favoriteID string,
	position int,
) error {
	return tx.Model(&e.FavoriteItem{}).
		Where("favorite_id = ? AND user_id = ?", favoriteID, uid).
		Update("position", position).
		Error
}
