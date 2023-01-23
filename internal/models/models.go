package models

type Models struct {
	Ingredients IngredientModel
	Recipes     RecipeModel
	Favorites   FavoriteModel
}

func New() Models {
	return Models{
		Ingredients: IngredientModel{Ingredients: Ingredients},
		Recipes:     RecipeModel{Recipes: Recipes},
		Favorites:   FavoriteModel{Favorites: Favorites},
	}
}
