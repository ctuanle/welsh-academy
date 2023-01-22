package models

import "time"

type RecipeIngredient struct {
	ID     int     `json:"id"`
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"` // g | ml | ...
}

type Recipe struct {
	ID          int                `json:"id"`
	Creator     int                `json:"creator"`
	Name        string             `json:"name"`
	Ingredients []RecipeIngredient `json:"ingredients"`
	Description string             `json:"string"`
	Created     time.Time          `json:"created"`
}

type RecipeModel struct {
	Recipes []Recipe
}

var Recipes = []Recipe{
	{
		ID:      1,
		Creator: 1,
		Name:    "Petits sablés",
		Ingredients: []RecipeIngredient{
			{
				ID:     1,
				Amount: 100,
				Unit:   "g",
			},
			{
				ID:     2,
				Amount: 150,
				Unit:   "g",
			},
			{
				ID:     4,
				Amount: 80,
				Unit:   "ml",
			},
		},
		Description: "This is a simple description",
		Created:     time.Now(),
	},
}

// GetAll() returns all existing recipes
func (m *RecipeModel) GetAll() ([]Recipe, error) {
	return m.Recipes, nil
}

// Insert() inserts new recipes
// and return this newly created recipe
func (m *RecipeModel) Insert(name, description string, creator int, ingredients []RecipeIngredient) (Recipe, error) {
	newRecipe := Recipe{
		ID:          len(m.Recipes) + 1,
		Name:        name,
		Description: description,
		Creator:     creator,
		Created:     time.Now(),
		Ingredients: ingredients,
	}

	m.Recipes = append(m.Recipes, newRecipe)
	return newRecipe, nil
}
