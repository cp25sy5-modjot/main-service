package favsvc

import (
	"errors"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	favrepo "github.com/cp25sy5-modjot/main-service/internal/favorite_item/repository"
	mapper "github.com/cp25sy5-modjot/main-service/internal/mapper"
	"github.com/google/uuid"
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
	repo *favrepo.Repository
}

func NewService(repo *favrepo.Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(input *m.FavoriteItemCreateInput) (*e.FavoriteItem, error) {
	maxPos, err := s.repo.GetMaxPosition(input.UserID)
	if err != nil {
		return nil, err
	}

	FavID := uuid.New().String()
	u := mapper.BuildFavObjectToCreate(FavID, input)
	u.Position = maxPos + 1

	favCreated, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}
	return favCreated, nil
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
	fav, err := s.repo.FindByID(uid, favID)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(uid, favID); err != nil {
		return err
	}

	return s.repo.ShiftLeftAfter(uid, fav.Position)
}

func (s *service) ReOrder(req *m.FavoriteItemReOrderInput) error {
	if len(req.ReOrderList) == 0 {
		return nil
	}

	seen := map[int]bool{}
	for _, item := range req.ReOrderList {
		if seen[item.Position] {
			return errors.New("duplicate position")
		}
		seen[item.Position] = true
	}

	for _, item := range req.ReOrderList {
		if item.FavoriteID == "" {
			continue
		}
		if item.Position <= 0 {
			continue
		}

		if err := s.repo.UpdatePosition(
			req.UserID,
			item.FavoriteID,
			item.Position,
		); err != nil {
			return err
		}
	}

	return nil
}
