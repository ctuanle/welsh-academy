package models

import (
	"database/sql"
)

type Models struct {
	Ingredients interface {
		GetAll() ([]*Ingredient, error)
		GetById(id int) (*Ingredient, error)
		Insert(ingredient *Ingredient) error
	}
	Recipes interface {
		GetAll(include map[int]struct{}, exclude map[int]struct{}) ([]*Recipe, error)
		Insert(recipe *Recipe) error
	}
	Favorites interface {
		GetAll(user_id int) ([]*Favorite, error)
		Insert(fav *Favorite) error
		Remove(favoriteId int) error
	}
}

// New() return models for app
func New(db *sql.DB) Models {
	return Models{
		Ingredients: IngredientModel{DB: db},
		Recipes:     RecipeModel{DB: db},
		Favorites:   FavoriteModel{DB: db},
	}
}
