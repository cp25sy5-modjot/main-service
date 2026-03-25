package favrepo

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(favorite *e.FavoriteItem) (*e.FavoriteItem, error)
	CreateTx(tx *gorm.DB, favorite *e.FavoriteItem) (*e.FavoriteItem, error)

	FindAll(uid string) ([]*e.FavoriteItem, error)
	FindByID(uid, favoriteID string) (*e.FavoriteItem, error)
	FindByIDTx(tx *gorm.DB, uid, favoriteID string) (*e.FavoriteItem, error)

	Update(favorite *e.FavoriteItem) (*e.FavoriteItem, error)

	Delete(uid, favoriteID string) error
	DeleteTx(tx *gorm.DB, uid, favoriteID string) error

	GetMaxPosition(uid string) (int, error)
	GetMaxPositionTx(tx *gorm.DB, uid string) (int, error)

	ShiftLeftAfter(uid string, pos int) error
	ShiftLeftAfterTx(tx *gorm.DB, uid string, pos int) error

	UpdatePosition(uid string, favoriteID string, position int) error
	UpdatePositionTx(tx *gorm.DB, uid string, favoriteID string, position int) error

	ResequencePositionsTx(tx *gorm.DB, uid string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(favorite *e.FavoriteItem) (*e.FavoriteItem, error) {
	if err := r.db.Create(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *repository) CreateTx(
	tx *gorm.DB,
	favorite *e.FavoriteItem,
) (*e.FavoriteItem, error) {
	if err := tx.Create(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *repository) FindAll(uid string) ([]*e.FavoriteItem, error) {
	var favorites []*e.FavoriteItem
	err := r.db.
		Where("user_id = ?", uid).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("category_id", "icon", "color_code", "category_name")
		}).
		Order("position ASC").
		Find(&favorites).Error
	return favorites, err
}

func (r *repository) FindByID(uid, favoriteID string) (*e.FavoriteItem, error) {
	var favorite e.FavoriteItem
	err := r.db.
		Where("favorite_id = ? AND user_id = ?", favoriteID, uid).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("category_id", "icon", "color_code", "category_name")
		}).
		First(&favorite).Error
	return &favorite, err
}

func (r *repository) FindByIDTx(
	tx *gorm.DB,
	uid,
	favoriteID string,
) (*e.FavoriteItem, error) {
	var favorite e.FavoriteItem
	err := tx.
		Where("favorite_id = ? AND user_id = ?", favoriteID, uid).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("category_id", "icon", "color_code", "category_name")
		}).
		First(&favorite).Error
	return &favorite, err
}

func (r *repository) Update(favorite *e.FavoriteItem) (*e.FavoriteItem, error) {
	if err := r.db.Save(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *repository) Delete(uid, favoriteID string) error {
	return r.db.
		Delete(&e.FavoriteItem{}, "favorite_id = ? AND user_id = ?", favoriteID, uid).
		Error
}

func (r *repository) DeleteTx(
	tx *gorm.DB,
	uid,
	favoriteID string,
) error {
	return tx.
		Delete(&e.FavoriteItem{}, "favorite_id = ? AND user_id = ?", favoriteID, uid).
		Error
}

func (r *repository) GetMaxPosition(uid string) (int, error) {
	var max int
	err := r.db.Model(&e.FavoriteItem{}).
		Where("user_id = ?", uid).
		Select("COALESCE(MAX(position), 0)").
		Scan(&max).Error
	return max, err
}

func (r *repository) GetMaxPositionTx(
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

func (r *repository) ShiftLeftAfter(uid string, pos int) error {
	return r.db.Model(&e.FavoriteItem{}).
		Where("user_id = ? AND position > ?", uid, pos).
		Update("position", gorm.Expr("position - 1")).
		Error
}

func (r *repository) ShiftLeftAfterTx(
	tx *gorm.DB,
	uid string,
	pos int,
) error {
	return tx.Model(&e.FavoriteItem{}).
		Where("user_id = ? AND position > ?", uid, pos).
		Update("position", gorm.Expr("position - 1")).
		Error
}

func (r *repository) UpdatePosition(
	uid string,
	favoriteID string,
	position int,
) error {
	return r.db.Model(&e.FavoriteItem{}).
		Where("favorite_id = ? AND user_id = ?", favoriteID, uid).
		Update("position", position).
		Error
}

func (r *repository) UpdatePositionTx(
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

func (r *repository) ResequencePositionsTx(tx *gorm.DB, uid string) error {
	return tx.Exec(`
		WITH ranked AS (
			SELECT favorite_id,
				   ROW_NUMBER() OVER (ORDER BY position) AS new_pos
			FROM favorite_items
			WHERE user_id = ?
		)
		UPDATE favorite_items f
		SET position = r.new_pos
		FROM ranked r
		WHERE f.favorite_id = r.favorite_id
	`, uid).Error
}
