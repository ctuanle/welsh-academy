package mocks

import (
	"database/sql"

	"ctuanle.ovh/welsh-academy/internal/models"
)

type MockFavoriteModel struct{}

func (m MockFavoriteModel) GetAll(user_id int) ([]*models.Favorite, error) {
	res := []*models.Favorite{}

	for _, f := range MockedFavorites {
		if f != nil && f.UserId == user_id {
			res = append(res, f)
		}
	}

	return res, nil
}

func (m MockFavoriteModel) Insert(fav *models.Favorite) error {
	fav.ID = len(MockedFavorites) + 1
	return nil
}

func (m MockFavoriteModel) Remove(favoriteId int) error {
	if favoriteId < 1 || favoriteId > len(MockedFavorites) || MockedFavorites[favoriteId] == nil {
		return sql.ErrNoRows
	}

	MockedFavorites[favoriteId] = nil
	return nil
}
