package favsvc

import (
	"errors"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	favrepo "github.com/cp25sy5-modjot/main-service/internal/favorite_item/repository"
	mapper "github.com/cp25sy5-modjot/main-service/internal/mapper"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(input *m.FavoriteItemCreateInput) (*e.FavoriteItem, error)

	GetAll(uid string) ([]*e.FavoriteItem, error)
	GetByID(uid string, fav_id string) (*e.FavoriteItem, error)
	Update(input *m.FavoriteItemUpdateInput) (*e.FavoriteItem, error)
	Delete(uid string, fav_id string) error
	ReOrder(req *m.FavoriteItemReOrderInput) error
}

type service struct {
	db *gorm.DB
	repo *favrepo.Repository
}

func NewService(db *gorm.DB, repo *favrepo.Repository) Service {
	return &service{db: db, repo: repo}
}

func (s *service) Create(input *m.FavoriteItemCreateInput) (*e.FavoriteItem, error) {
	var favCreated *e.FavoriteItem

	err := s.db.Transaction(func(tx *gorm.DB) error {
		maxPos, err := s.repo.GetMaxPositionTx(tx, input.UserID)
		if err != nil {
			return err
		}

		FavID := uuid.New().String()
		u := mapper.BuildFavObjectToCreate(FavID, input)
		u.Position = maxPos + 1

		favCreated, err = s.repo.CreateTx(tx, u)
		return err
	})

	return favCreated, err
}

func (s *service) GetAll(uid string) ([]*e.FavoriteItem, error) {
	return s.repo.FindAll(uid)
}

func (s *service) GetByID(uid string, fav_id string) (*e.FavoriteItem, error) {
	return s.repo.FindByID(uid, fav_id)
}

func (s *service) Update(input *m.FavoriteItemUpdateInput) (*e.FavoriteItem, error) {
	fav, err := s.repo.FindByID(input.UserID, input.FavoriteID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		fav.Title = *input.Title
	}
	if input.CategoryID != nil {
		fav.CategoryID = *input.CategoryID
	}
	if input.Price != nil {
		fav.Price = *input.Price
	}

	updatedFav, err := s.repo.Update(fav)
	if err != nil {
		return nil, err
	}
	return updatedFav, nil
}

func (s *service) Delete(uid, favID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		fav, err := s.repo.FindByIDTx(tx, uid, favID)
		if err != nil {
			return err
		}

		if err := s.repo.DeleteTx(tx, uid, favID); err != nil {
			return err
		}

		return s.repo.ShiftLeftAfterTx(tx, uid, fav.Position)
	})
}
func (s *service) ReOrder(req *m.FavoriteItemReOrderInput) error {
	if len(req.ReOrderList) == 0 {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {

		// 1) validate duplicate position ใน request
		seen := map[int]bool{}
		for _, item := range req.ReOrderList {
			if seen[item.Position] {
				return errors.New("duplicate position in request")
			}
			seen[item.Position] = true
		}

		// 2) move all affected rows to temp positions (negative)
		for _, item := range req.ReOrderList {
			tmpPos := -item.Position
			if err := s.repo.UpdatePositionTx(
				tx,
				req.UserID,
				item.FavoriteID,
				tmpPos,
			); err != nil {
				return err
			}
		}

		// 3) set real positions
		for _, item := range req.ReOrderList {
			if err := s.repo.UpdatePositionTx(
				tx,
				req.UserID,
				item.FavoriteID,
				item.Position,
			); err != nil {
				return err
			}
		}

		return nil
	})
}
