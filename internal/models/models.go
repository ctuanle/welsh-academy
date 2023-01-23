package models

import "database/sql"

type Models struct {
	Ingredients IngredientModel
	Recipes     RecipeModel
	Favorites   FavoriteModel
}

func New(db *sql.DB) Models {
	return Models{
		Ingredients: IngredientModel{DB: db},
		Recipes:     RecipeModel{DB: db},
		Favorites:   FavoriteModel{DB: db},
	}
}
