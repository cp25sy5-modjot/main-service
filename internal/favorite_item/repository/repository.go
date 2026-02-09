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

func (r *Repository) FindAll(uid string) ([]*e.FavoriteItem, error) {
	var favorites []*e.FavoriteItem
	err := r.db.Where("user_id = ?", uid).Find(&favorites).Error
	return favorites, err
}

func (r *Repository) FindByID(uid string, favorite_id string) (*e.FavoriteItem, error) {
	var favorite e.FavoriteItem
	err := r.db.Where("favorite_id = ? AND user_id = ?", favorite_id, uid).First(&favorite).Error
	return &favorite, err
}

func (r *Repository) Update(favorite *e.FavoriteItem) (*e.FavoriteItem, error) {
	if err := r.db.Save(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func (r *Repository) Delete(uid string, favorite_id string) error {
	return r.db.Delete(&e.FavoriteItem{}, "favorite_id = ? AND user_id = ?", favorite_id, uid).Error
}

func (r *Repository) GetMaxPosition(uid string) (int, error) {
	var max int
	err := r.db.Model(&e.FavoriteItem{}).
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

func (r *Repository) UpdatePosition(
	uid string,
	favID string,
	position int,
) error {
	return r.db.Model(&e.FavoriteItem{}).
		Where("favorite_id = ? AND user_id = ?", favID, uid).
		Update("position", position).
		Error
}
