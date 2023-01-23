package models

import (
	"database/sql"
)

type Favorite struct {
	ID       int `json:"id"`
	RecipeId int `json:"recipe_id"`
	UserId   int `json:"user_id"`
}

type FavoriteModel struct {
	DB *sql.DB
}

// GetAll() returns all favorite recipes of an user
func (m *FavoriteModel) GetAll(user_id int) ([]*Favorite, error) {
	query := `
		SELECT id, recipe_id, user_id
		FROM favorites
		WHERE user_id = $1
	`

	rows, err := m.DB.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favorites := []*Favorite{}

	for rows.Next() {
		var fav Favorite
		err := rows.Scan(&fav.ID, &fav.RecipeId, &fav.UserId)

		if err != nil {
			return nil, err
		}

		favorites = append(favorites, &fav)
	}

	return favorites, nil
}

// Insert() flags recipe_id as user_id favorite
func (m *FavoriteModel) Insert(fav *Favorite) error {
	query := `
		INSERT INTO favorites (recipe_id, user_id)
		VALUES ($1, $2)
		RETURNING id
	`

	return m.DB.QueryRow(query, fav.RecipeId, fav.UserId).Scan(&fav.ID)
}

// Remove() unflags recipe_id from user_id favorite
func (m *FavoriteModel) Remove(favoriteId int) error {
	if favoriteId < 1 {
		return sql.ErrNoRows
	}

	query := "DELETE FROM favorites WHERE id = $1"

	result, err := m.DB.Exec(query, favoriteId)
	if err != nil {
		return err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsEffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
